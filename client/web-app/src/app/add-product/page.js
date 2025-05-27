import React from "react";
import NavigationBar from "../Components/NavigationBar";
import EditableComponent from "../Components/Editable/EditableComponent";

export default function AddProduct()
{
    return(
        <>
        <NavigationBar />
        <EditableComponent edit={false} />
        </>
    )
}