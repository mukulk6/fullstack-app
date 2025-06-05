"use client";
import React, { useEffect, useState } from "react";
import { Card, CardContent, Grid, Typography, Box, Button, CircularProgress,Modal,Snackbar,Alert, Select,FormControl,MenuItem } from "@mui/material";
import AddIcon from '@mui/icons-material/Add';
import { useParams, useRouter } from "next/navigation";
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import NavigationBar from "../Components/NavigationBar";

const style = {
  position: 'absolute',
  top: '50%',
  left: '50%',
  transform: 'translate(-50%, -50%)',
  width: 400,
  bgcolor: 'background.paper',
  border: '2px solid #000',
  boxShadow: 24,
  p: 4,
};

const stylingOptions = {
	"outline":"none",
	"padding":"5px",
	"border":"1px solid lightgray"
}

export default function ProductsPage() {
	const [loading, setLoading] = useState(false)
	const [role, setRole] = useState('')
	const router = useRouter();
	 const [open, setOpen] = React.useState(false);
  const handleOpen = (id) => {setOpen(true),setDeleteId(id)}
  const handleClose = () => setOpen(false);
	const [openToastMessage,setOpenToastMessage] = useState(false);
	const [message,setMessage] = useState('')
	const fetchProducts = () => {
		 const token = localStorage.getItem("token");
		setLoading(true)
		fetch('http://localhost:8080/products', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json',
				 Authorization: `Bearer ${token}`,
			},
		})
			.then((res) => res.json()) // ✅ return the parsed JSON
			.then((data) => {
				setLoading(false);
				setData(data.products); 
				// ✅ save it to state
			})
			.catch((err) => console.log(err));
	};
	const [data, setData] = useState([]);
	const [deleteId,setDeleteId] = useState(null)
	const [options,setOptions] = useState([])

	const onClickAddProduct = () => {
		router.push('/add-product')
	}


	const onClickDelete = () => {
			const token = localStorage.getItem("token");
		    fetch(`http://localhost:8080/products/${deleteId}`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
			'Authorization':`Bearer ${token}`
        }
    })
    .then((res) => {
        if (!res.ok) {
            throw new Error('Failed to delete product');
        }
        return res.json();
    })
    .then((data) => {
				fetchProducts();
        setOpen(false);
				setOpenToastMessage(true);
				setMessage(data.message)
    })
    .catch((err) => {
        console.error("Error deleting product:", err);
    });
	}

		const Toast = ({ open, message, onClose }) => (
			<Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'right' }} open={open} autoHideDuration={3000} onClose={onClose}>
				<Alert onClose={onClose} severity={"success"} variant="filled" sx={{ width: '100%' }}>
					{message}
				</Alert>
			</Snackbar>
		);

		const handleCloseToastMessages = () => {
			setOpenToastMessage(false)
		}

		const onChangeUser = (val) => {
			console.log(val.target.value,'check value')
			const token = localStorage.getItem("token")
			fetch(`http://localhost:8080/products?admin=${encodeURIComponent(val.target.value)}`,{
				method:'GET',
				headers:{
					  'Content-Type': 'application/json',
						'Authorization':`Bearer ${token}`
				},
			}).then((res)=>res.json()).then((data)=>setData(data?.products)).catch((err)=>console.log(err))
		}

		const onClickCard = (id) => {
			router.push(`new-product/${id}`)
		}

		const getProductsList = (val) => {
			const token = localStorage.getItem("token")
			if(val !== '')
			{
				fetch(`http://localhost:8080/search?query=${val}`,{
					method:'GET',
						headers:{
					  'Content-Type': 'application/json',
				},
				}).then((res)=>res.json()).then((data)=>{
					if(data)
					{
						if(data.results && data.results.length > 0 && val !== "")
						{
							setData(data.results)
							if(val === "" || val.length ===0)
							{
								fetchProducts()
							}
						}
					}
				}).catch((err)=>console.log(err))
			}
		}

	useEffect(() => {
		var token = null;
		if(localStorage.getItem("token") && localStorage.getItem("token") !== '')
		{
			 token = localStorage.getItem("token")
		}
		fetchProducts();
		if (localStorage.getItem("role") && localStorage.getItem("role") !== "") {
			setRole(localStorage.getItem("role"))
		}
		if(localStorage.getItem("role") && localStorage.getItem("role") === 'Admin')
		{
			fetch('http://localhost:8080/admins/list',{
				    method: 'GET',
      headers: {
        'Content-Type': 'application/json',
				'Authorization':`Bearer ${token}`
      },
			}).then((res)=>res.json()).then((data)=>setOptions(data)).catch((err)=>console.log(err))
		}
	}, [])
	return (
		<>
			<NavigationBar />
						{role === 'Admin' &&<Box sx={{ display: 'flex', justifyContent: 'end',width:'90%',marginTop:'12px' }}>
							<Grid container>
						<Grid item>
							<FormControl>
								{/* <Select
								labelId="demo-simple-select-label"
								id="demo-simple-select"
								label="Age"
								value={options}
							>
								{options.length > 0 && options.map((option) => (
									<MenuItem key={option.value} value={option.value}>
										{option.label}
									</MenuItem>
								))}
								</Select> */}
								<select aria-placeholder="WHat" onChange={(e)=>onChangeUser(e)} style={stylingOptions} >
									<option value={''}>All</option>
									{options.length > 0  && options.map((obj,ind)=>{
									return(<option value={obj.label} key={ind}>{obj.label}</option>)
								})}</select>
						</FormControl>
						 
							<>
							<input placeholder="Search for Products" style={stylingOptions} type="text" onChange={(e)=>getProductsList(e.target.value)} />
							</>
						</Grid>
						</Grid>
					</Box>}
			<Grid container spacing={2} paddingTop={3} justifyContent={'center'}>
		
				{loading && <Box sx={{ display: 'flex', justifyContent: 'center' }}><CircularProgress /></Box>}
				{data && loading === false && data.length > 0 && data.map((obj, ind) => {
					return (
						<>
						<Grid item key={ind}>
							<Card sx={{
								height: 300,
								width: 300,
								mb: 2,
								cursor: 'pointer',
								boxShadow: '0 4px 20px rgba(0,0,0,0.1)',
								transition: 'transform 0.3s ease, box-shadow 0.3s ease',
								'&:hover': {
									transform: 'scale(1.05)',
									boxShadow: '0 8px 30px rgba(0,0,0,0.2)',
								},
							}} onClick={()=>role === 'User' ? onClickCard(obj.id) : ""} key={ind} variant="outlined"
							><CardContent sx={{ height: '100%' }}><Typography gutterBottom sx={{ color: 'text.secondary', fontSize: 14 }}>
								{obj.name}
							</Typography> <Typography variant="h6" component="div">
										{obj.description}</Typography>
									<Typography sx={{ color: 'text.secondary', }}>{obj.price + " " + "INR"}</Typography>
									<Typography sx={{ color: 'text.secondary', mb: 1.2 }}>Quantity:{" "}{obj.quantity}</Typography>
									<div style={{ textAlign: 'end' }}>
									{role === 'Admin' && localStorage.getItem("role") && localStorage.getItem("role") === 'Admin' && (<><Typography sx={{ color: 'text.secondary', mb: 1.2 }}>Product Id:{" "}{obj.id}</Typography>
									<div style={{ textAlign: 'end' }}></div></>)}
										{role === 'Admin' && localStorage.getItem("role") && localStorage.getItem("role") === 'Admin' && <Button onClick={() => router.push(`/new-product/${obj.id}`)} style={{ textTransform: 'none', fontSize: '14px', marginRight: '7px' }} endIcon={<EditIcon />} variant="contained">Edit</Button>}
										{role === 'Admin' && localStorage.getItem("role") && localStorage.getItem("role") === 'Admin' && <Button onClick={()=>handleOpen(obj.id)} endIcon={<DeleteIcon />
									} style={{ backgroundColor: 'red', textTransform: 'none', fontSize: '14px' }} variant="contained">Delete</Button>}
									</div>
								</CardContent></Card></Grid>
								</>
								)
				})}
				<Modal
					open={open}
					onClose={handleClose}
					aria-labelledby="modal-modal-title"
					aria-describedby="modal-modal-description"
				>
					<Box sx={style}>
						<Typography id="modal-modal-title" variant="h6" component="h4">
							Are you sure you want to delete the product details?
						</Typography>
						<div style={{textAlign:'end'}}>
						<Button onClick={()=>onClickDelete()} style={{marginRight:'7px'}} variant="contained">Yes</Button>
						<Button variant="outlined" onClick={()=>handleClose()}>No</Button>
						</div>
					</Box>
				</Modal>
				{openToastMessage && <Toast message={message} open={openToastMessage} onClose={handleCloseToastMessages} />}
				{
					role === 'Admin' && loading === false && (
						<>
							<Grid item>
								<Card sx={{
									height: 300,
									width: 300
								}}>
									<CardContent style={{ height: '100%' }}>
										<div onClick={() => onClickAddProduct()} style={{
											display: 'flex',
											justifyContent: 'center',
											alignItems: 'center'
										}}><AddIcon sx={{ cursor: 'pointer', fontSize: '40px', borderRadius: '50%', backgroundColor: 'lightgray' }} />
										</div>
										<div style={{
											display: 'flex',
											justifyContent: 'center',
											fontWeight: 400,
											alignItems: 'center', paddingTop: '5px', color: 'rgba(0, 0, 0, 0.6)'
										}} >Add Product</div>
									</CardContent>
								</Card>
							</Grid>
						</>
					)
				}
			</Grid>
		</>
	)
}