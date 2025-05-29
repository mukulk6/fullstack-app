"use client";
import React, { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { Grid, Box, AppBar, Toolbar, IconButton, Button, Card, Typography, CardContent, Snackbar, Alert } from "@mui/material";
import ArrowBackRoundedIcon from '@mui/icons-material/ArrowBackRounded';

export default function EditableComponent({ edit, id, data }) {
	const [open, setOpen] = useState(false);
	const [message, setMessage] = useState('');
	const router = useRouter();
	const [productId,setProductId] = useState(id)

	const onClickRedirect = () => {
		router.push('/products-home')
	}

	const handleClose = () => {
		setOpen(!open)
	}
	// name, description, price, quantity
	const [inputValues, setInputValues] = useState({
		"name": '',
		"description": '',
		"price": '',
		"quantity": ''
	})

	const inputBoxStyles = {
		"outline": "none",
		"padding": "3px"
	}

	const onChangeValues = (e) => {
		setInputValues({ ...inputValues, [e.target.name]: e.target.value })
	}
	const onClickSubmitValues = () => {
			const token = localStorage.getItem("token")
				let newObj = {
			"name": inputValues.name,
			"description": inputValues.description,
			"price": parseFloat(inputValues.price),
			"quantity": parseFloat(inputValues.quantity)
		}
		if(edit === false)
		{
		fetch('http://localhost:8080/products', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
				'Authorization': `Bearer ${token}`
			},
			body: JSON.stringify(newObj),
		}).then((res) => res.json()).then((data) => {
			setOpen(true), setMessage(data.message),
				setInputValues({
					name: '',
					description: '',
					price: '',
					quantity: ''
				});
		}).catch((err) => console.log(err))
	}
if (edit) {
	fetch(`http://localhost:8080/products/${productId}`, {
		method: 'PUT',
		headers: {
			'Content-Type': 'application/json',
			'Authorization':token
		},
		body: JSON.stringify(newObj),
	})
		.then((res) => {
			if (!res.ok) {
				throw new Error("Failed to update the product.");
			}
			return res.json();
		})
		.then((data) => {
			setOpen(true);
			setMessage(data.message);

			// Optional: only reset if needed
			setInputValues({
				name: '',
				description: '',
				price: '',
				quantity: ''
			});
		})
		.catch((err) => {
			console.error("Error updating product:", err);
		});
}
	}

	const Toast = ({ open, message, onClose }) => (
		<Snackbar anchorOrigin={{ vertical: 'top', horizontal: 'right' }} open={open} autoHideDuration={3000} onClose={onClose}>
			<Alert onClose={onClose} severity={"success"} variant="filled" sx={{ width: '100%' }}>
				{message}
			</Alert>
		</Snackbar>
	);

	useEffect(() => {
		if (data && Object.keys(data).length > 0) {
			setInputValues({
				name: data.name || '',
				description: data.description || '',
				price: data.price || '',
				quantity: data.quantity || ''
			});
		}
		if(id)
		{
			setProductId(id)
		}
	}, [data]);
	return (
		<>
			<div>
				<Grid container justifyContent={'center'} paddingTop={2}>
					<Grid item>
						<div style={{ display: 'flex', alignItems: 'center' }}>
							<ArrowBackRoundedIcon style={{ marginRight: '5px', cursor: 'pointer', color: '#1976d2', fontSize: '30px' }} /><span style={{ fontSize: '24px', fontWeight: 'bold', cursor: 'pointer' }} onClick={() => onClickRedirect()}>Back</span>
						</div>
						<Card sx={{ width: '400px', height: '400px', marginTop: '5px' }}>
							<CardContent>
								<Typography style={{ color: 'GrayText', fontWeight: 'bold', textAlign: 'center', fontSize: '20px', paddingBottom: '5px' }}>{edit ? "Edit" :"Add"} Product Details</Typography>
								<Grid paddingBottom={1} container justifyContent={'center'} alignItems={"center"} columnSpacing={1}>
									<Grid item size={{ lg: 3 }}><span>Name:</span></Grid>
									<Grid item lg={4}><input onChange={(e) => onChangeValues(e)} value={inputValues.name} name="name" type="text" style={inputBoxStyles} /></Grid>
								</Grid>
								<Grid paddingBottom={1} container justifyContent={'center'} alignItems={"center"} columnSpacing={1}>
									<Grid item size={{ lg: 3 }}><span>Description:</span></Grid>
									<Grid item lg={4}><input onChange={(e) => onChangeValues(e)} value={inputValues.description} name="description" type="text" style={inputBoxStyles} /></Grid>
								</Grid>
								<Grid paddingBottom={1} container justifyContent={'center'} alignItems={"center"} columnSpacing={1}>
									<Grid item size={{ lg: 3 }}><span>Price:</span></Grid>
									<Grid item lg={4}><input min={0} onChange={(e) => onChangeValues(e)} value={inputValues.price} name="price" type="number" style={inputBoxStyles} /></Grid>
								</Grid>
								<Grid paddingBottom={1} container justifyContent={'center'} alignItems={"center"} columnSpacing={1}>
									<Grid item size={{ lg: 3 }}><span>Quantity:</span></Grid>
									<Grid item lg={4}><input onChange={(e) => onChangeValues(e)} min={0} value={inputValues.quantity} name="quantity" type="number" style={inputBoxStyles} /></Grid>
								</Grid>
								<Grid container justifyContent={"end"}>
									<Grid item lg={12}>
										<Button onClick={() => onClickSubmitValues()} variant="outlined">{edit ? "Edit" : "Submit"}</Button>
									</Grid>
								</Grid>
								{open && <Toast open={open} message={message} onClose={handleClose} />}
							</CardContent>
						</Card>
					</Grid>
				</Grid>
			</div>
		</>
	)
}