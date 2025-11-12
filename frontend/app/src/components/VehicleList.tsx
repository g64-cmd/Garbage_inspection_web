import React, { useState, useEffect } from 'react';
import { getVehicles } from '../services/api';
import { Link as RouterLink } from 'react-router-dom';
import {
  List,
  ListItem,
  ListItemText,
  CircularProgress,
  Alert,
  Paper,
  ListItemIcon,
  Typography,
  Box,
  ListItemButton,
} from '@mui/material';
import LocalShippingIcon from '@mui/icons-material/LocalShipping';
import BatteryStdIcon from '@mui/icons-material/BatteryStd';
import PowerIcon from '@mui/icons-material/Power';

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

// --- Helper Component for Status ---

const VehicleStatusInfo: React.FC<{ status: VehicleStatus | null }> = ({ status }) => {
  if (!status) {
    return <Typography variant="body2" color="text.secondary">Status: Unknown</Typography>;
  }

  return (
    <Box sx={{ display: 'flex', alignItems: 'center', gap: 2, mt: 1 }}>
      <Box sx={{ display: 'flex', alignItems: 'center' }}>
        <BatteryStdIcon fontSize="small" sx={{ mr: 0.5 }} />
        <Typography variant="body2" color="text.secondary">
          {status.battery.toFixed(1)}%
        </Typography>
      </Box>
      <Box sx={{ display: 'flex', alignItems: 'center' }}>
        <PowerIcon fontSize="small" sx={{ mr: 0.5 }} />
        <Typography variant="body2" color="text.secondary">
          {status.state}
        </Typography>
      </Box>
    </Box>
  );
};


// --- Main Component ---

const VehicleList: React.FC = () => {
  const [vehicles, setVehicles] = useState<Vehicle[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchVehicles = async () => {
      try {
        setLoading(true);
        const data = await getVehicles();
        setVehicles(data || []);
      } catch (err) {
        setError('Failed to fetch vehicles. Please try again later.');
        console.error(err);
      } finally {
        setLoading(false);
      }
    };

    fetchVehicles();
  }, []);

  if (loading) {
    return <CircularProgress />;
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Paper elevation={3}>
      <List>
        {vehicles.length > 0 ? (
          vehicles.map((vehicle) => (
            <ListItemButton key={vehicle.id} component={RouterLink} to={`/vehicles/${vehicle.id}`} divider>
              <ListItemIcon>
                <LocalShippingIcon fontSize="large" />
              </ListItemIcon>
              <ListItemText
                primary={vehicle.name}
                secondary={
                  <>
                    <Typography component="span" variant="body2" color="text.primary">
                      Model: {vehicle.model}
                    </Typography>
                    <VehicleStatusInfo status={vehicle.current_status} />
                  </>
                }
              />
            </ListItemButton>
          ))
        ) : (
          <ListItem>
            <ListItemText primary="No vehicles found." />
          </ListItem>
        )}
      </List>
    </Paper>
  );
};

export default VehicleList;
