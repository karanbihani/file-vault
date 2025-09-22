import StatsDisplay from "./StatsDisplay";

const StatsPage = () => {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white">
        📊 Your Statistics
      </h1>
      <StatsDisplay />
    </div>
  );
};

export default StatsPage;
