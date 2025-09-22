// Enhanced rate limiter utility with retry mechanism
class RateLimiter {
  private queue: Array<() => Promise<any>> = [];
  private isProcessing = false;
  private lastRequestTime = 0;
  private readonly minInterval = 500; // 500ms = 2 requests per second
  private readonly maxRetries = 3;
  private readonly retryDelay = 1000; // 1 second base delay

  async execute<T>(request: () => Promise<T>, retryCount = 0): Promise<T> {
    return new Promise((resolve, reject) => {
      this.queue.push(async () => {
        try {
          const result = await request();
          resolve(result);
        } catch (error: any) {
          // Check if it's a rate limit error (429) or network error
          const isRateLimitError = error?.response?.status === 429;
          const isNetworkError = !error?.response && error?.code;

          if (
            (isRateLimitError || isNetworkError) &&
            retryCount < this.maxRetries
          ) {
            console.warn(
              `Rate limit/network error, retrying... (${retryCount + 1}/${
                this.maxRetries
              })`
            );

            // Exponential backoff: 1s, 2s, 4s
            const backoffDelay = this.retryDelay * Math.pow(2, retryCount);

            setTimeout(() => {
              this.execute(request, retryCount + 1)
                .then(resolve)
                .catch(reject);
            }, backoffDelay);
          } else {
            reject(error);
          }
        }
      });

      this.processQueue();
    });
  }

  private async processQueue() {
    if (this.isProcessing || this.queue.length === 0) {
      return;
    }

    this.isProcessing = true;

    while (this.queue.length > 0) {
      const now = Date.now();
      const timeSinceLastRequest = now - this.lastRequestTime;

      if (timeSinceLastRequest < this.minInterval) {
        await new Promise((resolve) =>
          setTimeout(resolve, this.minInterval - timeSinceLastRequest)
        );
      }

      const request = this.queue.shift();
      if (request) {
        this.lastRequestTime = Date.now();
        await request();
      }
    }

    this.isProcessing = false;
  }

  // Method to handle burst uploads with longer delays
  async executeUpload<T>(
    request: () => Promise<T>,
    uploadDelay = 1000
  ): Promise<T> {
    return new Promise((resolve, reject) => {
      this.queue.push(async () => {
        try {
          const result = await request();
          resolve(result);
          // Add extra delay after uploads
          await new Promise((r) => setTimeout(r, uploadDelay));
        } catch (error: any) {
          const isRateLimitError = error?.response?.status === 429;
          if (isRateLimitError) {
            console.warn("Upload rate limit hit, waiting longer...");
            await new Promise((r) => setTimeout(r, 2000)); // Wait 2 seconds
            // Retry once for uploads
            try {
              const retryResult = await request();
              resolve(retryResult);
            } catch (retryError) {
              reject(retryError);
            }
          } else {
            reject(error);
          }
        }
      });

      this.processQueue();
    });
  }
}

// Create a singleton instance
export const rateLimiter = new RateLimiter();

// Debounce utility for search inputs
export function debounce<T extends (...args: any[]) => any>(
  func: T,
  wait: number
): (...args: Parameters<T>) => void {
  let timeout: ReturnType<typeof setTimeout>;

  return (...args: Parameters<T>) => {
    clearTimeout(timeout);
    timeout = setTimeout(() => func(...args), wait);
  };
}
