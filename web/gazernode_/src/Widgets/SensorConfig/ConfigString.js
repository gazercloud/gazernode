import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import {Button} from "@material-ui/core";
import DialogLookup from "../DialogLookup";

export default function ConfigString(props) {
    return (
        <Grid container alignItems="center" spacing={1}>
            <Grid item style={{minWidth: "100px"}}>{props.Meta.display_name}:</Grid>
            <Grid item>
                <TextField value={props.Data} onChange={(ev)=>{
                    props.OnChangedValue(props.Meta.name, ev.target.value)
                }}/>
            </Grid>
            <Grid item>
                <DialogLookup Entity={props.Meta.format} OnSelected={(v)=>{
                    props.OnChangedValue(props.Meta.name, v)
                }}>
                </DialogLookup>
            </Grid>
        </Grid>
    )
}
