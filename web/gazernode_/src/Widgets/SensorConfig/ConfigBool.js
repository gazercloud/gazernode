import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import {Checkbox, Tooltip, Zoom} from "@material-ui/core";

export default function ConfigBool(props) {
    const handleChange = (event) => {
        props.OnChangedValue(props.Meta.name, event.target.checked)
    };

    const val = () => {
        if (props.Data === undefined)
            return false
        return props.Data
    }

    return (
        <Grid container alignItems="center" spacing={1}>
            <Grid item style={{minWidth: "100px"}}>
                <div>{props.Meta.display_name}:</div>
            </Grid>
            <Grid item>
                <Checkbox
                    checked={val()}
                    onChange={handleChange}
                >
                </Checkbox>
            </Grid>
        </Grid>
    )
}
