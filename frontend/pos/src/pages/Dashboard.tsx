import { useEffect, useState } from 'react';
import { useAuth } from '../context/AuthContext';
import axios from 'axios';
import { Container, Typography } from '@mui/material';
const Dashboard = () => {

  return (
    <>
      <Container sx={{ marginTop: 4 }}>
        <Typography variant="h5">Selamat datang di Dashboard!</Typography>
      </Container>
    </>
  );
};

export default Dashboard;