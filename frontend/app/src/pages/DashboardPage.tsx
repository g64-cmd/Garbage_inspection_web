import React, { useState, useEffect } from 'react';
import { Typography, Box, Grid, Paper, CircularProgress, Alert } from '@mui/material';
import VehicleList from '../components/VehicleList';
import { getAllDecisionLogs } from '../services/api';
import { Pie } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip,
  Legend,
  Title,
} from 'chart.js';

// Register Chart.js components
ChartJS.register(ArcElement, Tooltip, Legend, Title);

// --- Data Models ---
interface ServerDecision {
  action: string;
}

interface DecisionLog {
  server_decision: ServerDecision;
}

interface DecisionLogResponse {
  logs: DecisionLog[];
}

// --- Chart Component ---

const DecisionPieChart: React.FC = () => {
  const [chartData, setChartData] = useState<any>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchData = async () => {
      try {
        setLoading(true);
        const response: DecisionLogResponse = await getAllDecisionLogs();
        const logs = response.logs || [];
        
        // Process data for the chart
        const actionCounts = logs.reduce((acc, log) => {
          const action = log.server_decision.action;
          acc[action] = (acc[action] || 0) + 1;
          return acc;
        }, {} as { [key: string]: number });

        const labels = Object.keys(actionCounts);
        const data = Object.values(actionCounts);

        setChartData({
          labels,
          datasets: [
            {
              label: 'Decision Actions',
              data,
              backgroundColor: [
                'rgba(255, 99, 132, 0.7)',
                'rgba(54, 162, 235, 0.7)',
                'rgba(255, 206, 86, 0.7)',
                'rgba(75, 192, 192, 0.7)',
                'rgba(153, 102, 255, 0.7)',
                'rgba(255, 159, 64, 0.7)',
              ],
              borderColor: [
                'rgba(255, 99, 132, 1)',
                'rgba(54, 162, 235, 1)',
                'rgba(255, 206, 86, 1)',
                'rgba(75, 192, 192, 1)',
                'rgba(153, 102, 255, 1)',
                'rgba(255, 159, 64, 1)',
              ],
              borderWidth: 1,
            },
          ],
        });

      } catch (err) {
        setError('Failed to fetch decision logs for chart.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return <CircularProgress />;
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Paper elevation={3} sx={{ p: 3 }}>
      <Typography variant="h6" gutterBottom>
        Decision Statistics
      </Typography>
      <Box sx={{ height: '400px', position: 'relative' }}>
        {chartData && <Pie data={chartData} options={{ responsive: true, maintainAspectRatio: false }} />}
      </Box>
    </Paper>
  );
};


// --- Main Dashboard Page ---

const DashboardPage: React.FC = () => {
  return (
    <Box>
      <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(12, 1fr)', gap: 3 }}>
        {/* Vehicle List taking up 2/3 of the width */}
        <Box sx={{ gridColumn: 'span 8' }}>
          <Typography variant="h4" gutterBottom>
            Vehicles Overview
          </Typography>
          <VehicleList />
        </Box>

        {/* Chart taking up 1/3 of the width */}
        <Box sx={{ gridColumn: 'span 4' }}>
           <DecisionPieChart />
        </Box>
      </Box>
    </Box>
  );
};

export default DashboardPage;
