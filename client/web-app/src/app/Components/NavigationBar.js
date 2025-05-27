"use client";
import React from "react";
import { Box,AppBar,Toolbar,IconButton,Button } from "@mui/material";

export default function NavigationBar() {
    return (
        <>
            <div>
            <Box sx={{ flexGrow: 1 }}>
                <AppBar position="static">
                    <Toolbar>
                        <IconButton
                            size="large"
                            edge="start"
                            color="inherit"
                            aria-label="menu"
                            sx={{ mr: 2 }}
                        >
                        </IconButton>
                        <Button color="inherit">Logout</Button>
                    </Toolbar>
                </AppBar>
            </Box>
            </div>
        </>
    )
}