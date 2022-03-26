import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";

export default function ConfigText(props) {
    return (
        <Grid container alignItems="center" spacing={1}>
            <Grid item style={{minWidth: "100px"}}>{props.Meta.display_name}:</Grid>
            <Grid item>
                <TextField rows={10} style={{minWidth: "300px"}} multiline={true} value={props.Data} onChange={(ev)=>{
                    props.OnChangedValue(props.Meta.name, ev.target.value)
                }}/>
            </Grid>
        </Grid>
    )
}
