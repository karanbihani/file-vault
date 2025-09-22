import StatsDisplay from "./StatsDisplay";
import { FiBarChart } from "react-icons/fi";

const StatsPage = () => {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6 text-gray-800 dark:text-white flex items-center">
        <FiBarChart className="mr-3" />
        Your Statistics
      </h1>
      <StatsDisplay />
    </div>
  );
};

export default StatsPage;
