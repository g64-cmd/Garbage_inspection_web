import React, { useState, useEffect } from 'react';
import { useParams } from 'react-router-dom';
import { getVehicleById, getDecisionLogsByVehicleId } from '../services/api';
import {
  Typography,
  Box,
  CircularProgress,
  Alert,
  Grid,
  Card,
  CardMedia,
  CardContent,
  Paper,
} from '@mui/material';

// --- Data Models to match backend ---

interface Position {
  lat: number;
  lng: number;
}

interface VehicleStatus {
  timestamp: number;
  position: Position;
  battery: number;
  state: string;
}

interface Vehicle {
  id: string;
  name: string;
  model: string;
  current_status: VehicleStatus | null;
}

interface ServerDecision {
  image_id: string;
  action: string;
  confidence: number;
  reason: string;
}

interface DecisionLog {
  id: string;
  vehicle_id: string;
  timestamp: string;
  image_url: string;
  server_decision: ServerDecision;
}

interface DecisionLogResponse {
  logs: DecisionLog[];
  total: number;
}

// --- Main Component ---

const VehicleDetailPage: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [vehicle, setVehicle] = useState<Vehicle | null>(null);
  const [logs, setLogs] = useState<DecisionLog[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!id) return;

    const fetchData = async () => {
      try {
        setLoading(true);
        const vehicleData = await getVehicleById(id);
        setVehicle(vehicleData);

        const logsData: DecisionLogResponse = await getDecisionLogsByVehicleId(id);
        setLogs(logsData.logs || []);

      } catch (err) {
        setError('Failed to fetch vehicle details. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, [id]);

  if (loading) {
    return <CircularProgress />;
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  if (!vehicle) {
    return <Alert severity="info">Vehicle not found.</Alert>;
  }

  return (
    <Box>
      <Paper elevation={3} sx={{ p: 3, mb: 4 }}>
        <Typography variant="h4" gutterBottom>
          {vehicle.name}
        </Typography>
        <Typography variant="h6" color="text.secondary">
          Model: {vehicle.model}
        </Typography>
      </Paper>

      <Typography variant="h5" gutterBottom sx={{ mt: 4, mb: 2 }}>
        Decision Logs
      </Typography>

      {logs.length > 0 ? (
        <Box sx={{ display: 'grid', gridTemplateColumns: 'repeat(12, 1fr)', gap: 3 }}>
          {logs.map((log) => (
            <Box sx={{ gridColumn: { xs: 'span 12', sm: 'span 6', md: 'span 4' } }} key={log.id}>
              <Card>
                <CardMedia
                  component="img"
                  height="200"
                  image={log.image_url} // Directly use the image URL
                  alt={`Decision log ${log.id}`}
                />
                <CardContent>
                  <Typography gutterBottom variant="body2" color="text.secondary">
                    {new Date(log.timestamp).toLocaleString()}
                  </Typography>
                  <Typography variant="h6" component="div">
                    Action: {log.server_decision.action}
                  </Typography>
                  <Typography variant="body1" color="text.primary">
                    Confidence: {(log.server_decision.confidence * 100).toFixed(1)}%
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Reason: {log.server_decision.reason}
                  </Typography>
                </CardContent>
              </Card>
            </Box>
          ))}
        </Box>
      ) : (
        <Alert severity="info">No decision logs found for this vehicle.</Alert>
      )}
    </Box>
  );
};

export default VehicleDetailPage;
