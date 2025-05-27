import React from "react";
import NavigationBar from "../Components/NavigationBar";
import EditableComponent from "../Components/Editable/EditableComponent";
import { useParams } from "next/navigation";

export default function EditPage()
{

    return(
        <>
        <NavigationBar />
        <EditableComponent  />
        </>
    )
}