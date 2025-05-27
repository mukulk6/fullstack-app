'use client';
import React, { useState } from "react";
import styles from "./loginpage.module.css"
import { Button, Grid,Snackbar,Alert} from "@mui/material";
import { useRouter } from "next/navigation";



export default function LoginPage()
{
    const [clickSignUp, setClickSignUp] = React.useState(false);
    const [password, setPassword] = React.useState('');
    const [email,setEmail] = React.useState('');
    const [role,setRole] = React.useState('')
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    const [openSnackbar, setOpenSnackbar] = useState(false);
    const [message, setMessage] = useState('');
    const [error, setError] = useState(false);
    const router = useRouter();

    console.log(error);

    const validateEmail = (e) => {
        setEmail(e.target.value);
        return emailRegex.test(e.target.value)
    }

const handleClose = () => {
    setOpenSnackbar(false)
}



const Toast = ({ open, message, onClose,error }) => (
  <Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'right' }} open={open} autoHideDuration={3000} onClose={onClose}>
    <Alert onClose={onClose} severity={error ? "error" : "success"} variant="filled" sx={{ width: '100%' }}>
      {message}
    </Alert>
  </Snackbar>
);

    const signIn = () => {
        let newObj = {
            email: email,
            password: password
        }
        fetch('http://localhost:8080/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(newObj),
        }).then((res) => res.json())
            .then((data) => {
                setMessage(data.message)
                setOpenSnackbar(true);
                setError(data?.error);
                if (data.error !== undefined && data.error) {
                    setError(data?.error)
                }
                if(data.role !== undefined)
                {
                localStorage.setItem("role",data.role);
                // fetchProducts();
                router.push('/products-home')
            }
            if(data.token)
            {
                localStorage.setItem('token',data.token)
            }
            })
            .catch((err) => console.log(err, 'check error message'))
    }

    const signUp = () => {
        let newObj = {   
            email:email,
            password:password,
            role:role
        }
fetch('http://localhost:8080/signup', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
  },
  body: JSON.stringify(newObj),
})
  .then((res) => res.json())
  .then((data) => {
    setMessage(data.message);         // Store message in state
    setOpenSnackbar(true);
    setError(data?.error);
    if(data.error !== undefined && data.error)
        {
            setError(error)
        }           // Open Snackbar // Optional: call external toast handler
  })
  .catch((err) =>{console.log(err,'check error')})}
    return(
        <>
        <div >
        <Grid container justifyContent={'center'}>
            <Grid sx={{backgroundColor:'#fff'}} item size={3} height={'300px'} border={'1px solid lightgray'} padding={'12px'}>
              <h6 className={styles.titleHeader}>{clickSignUp ? 'Sign Up': 'Login'}</h6> 
              <div className={styles.inputBox}>
                <span>*</span><span style={{paddingRight:'5px',color:'#5c5c5c',fontWeight:700}}>Email:</span>
                </div>
                <div style={{paddingBottom:'5px'}}>
                <input value={email} onChange={(e)=>validateEmail(e)} placeholder="Enter your email" type="email" className={styles.inputElement} />
                </div>
                {email.length > 0 && emailRegex.test(email) === false && <span style={{color:'red',fontSize:'12px'}}>Enter a valid email.</span>}
                <div className={styles.inputBox}>
                <span>*</span><span style={{paddingRight:'5px',color:'#5c5c5c',fontWeight:700}}>Password:</span>
                </div>
                <div style={{paddingBottom:'10px'}}>
                <input value={password} onChange={(e)=>setPassword(e.target.value)} placeholder="Enter your password" type="password" className={styles.inputElement} />
                {password.length !== 0 && password.length < 8 && <span style={{color:'red',paddingTop:'5px',fontSize:'12px'}}>Minimum length of characters should be 8.</span>}
                </div>
                {
                    clickSignUp && <>
                    <span style={{paddingRight:'5px',color:'#5c5c5c',fontWeight:700}}>*Role:</span>
                    <div style={{paddingTop:'7px',paddingBottom:'7px'}}>
                    <select defaultValue={''} value={role} onChange={(e)=>setRole(e.target.value)} className={styles.selectComponent}><option>Admin</option><option>User</option></select>
                    </div>
                    </>
                }
                {openSnackbar && <Toast error={error} onClose={handleClose} open={openSnackbar} message={message} />}
                <div>
                <Button onClick={clickSignUp ? signUp : signIn} variant="contained">{clickSignUp ?"Sign Up":"Sign In"}</Button>
                </div>
             
                <div style={{fontSize:'14px',paddingTop:'5px'}}>
                If you're a new user please <span onClick={()=>setClickSignUp(!clickSignUp)} style={{textDecoration:'underline',cursor:'pointer',color:'blue'}}>click here</span>
                </div>             
            </Grid>
        </Grid>
        </div>
        </>
    )
}