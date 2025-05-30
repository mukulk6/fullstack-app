"use client";
import React, { useEffect,useState } from "react";
import { Box,AppBar,Toolbar,IconButton,Button } from "@mui/material";
import AccountCircle from '@mui/icons-material/AccountCircle';
import SupervisorAccountIcon from '@mui/icons-material/SupervisorAccount';
import { useRouter } from "next/navigation";
export default function NavigationBar() {

     const [email, setEmail] = useState('');
     const [userRole, setRole] = useState('');
     const router = useRouter();

     const logOutUser = () => {
        localStorage.removeItem("token");
        localStorage.removeItem("role");
        localStorage.removeItem("email");
        router.push('/home')


     }

    useEffect(()=>{
   if (typeof window !== 'undefined') {
      const storedEmail = localStorage.getItem('email');
      const storedRole = localStorage.getItem('role');
      if (storedEmail && storedEmail !== '') {
        setEmail(storedEmail);
      }
      if(storedRole && storedRole !== '')
      {
        setRole(storedRole)
      }
    }  
    },[])

    return (
    <>
    <div>
  <Box sx={{ flexGrow: 1 }}>
    <AppBar position="static">
      <Toolbar>
        {/* Left side content (if any) can go here */}

        {/* Spacer to push next items to the right */}
        <Box sx={{ flexGrow: 1 }} />

        {/* User info and Logout at the end */}
        <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
          <div style={{ textAlign: 'right' }}>
            <span style={{ fontSize: '14px' }}>Welcome</span>
            <br />
            <span style={{ fontSize: '13px' }}>{email}</span>
          </div>
          <IconButton
            size="large"
            aria-label="account of current user"
            color="inherit"
          >
            {userRole === 'Admin' ? <SupervisorAccountIcon />: <AccountCircle />}
          </IconButton>
          <Button color="inherit" onClick={()=>logOutUser()}>Logout</Button>
        </Box>
      </Toolbar>
    </AppBar>
  </Box>
  </div>
</>

    )
}