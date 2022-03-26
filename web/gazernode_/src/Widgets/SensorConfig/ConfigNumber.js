import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import {Tooltip, Zoom} from "@material-ui/core";

export default function ConfigNumber(props) {
    const [outOfRange, setOutOfRange] = useState(false)

    const color = () => {
        if (outOfRange) {
            return "#f00"
        } else {
            return undefined
        }
    }

    const tooltip = () => {
        if (outOfRange) {
            return "out of range. min:" + props.Meta.min_value + " max:" + props.Meta.max_value;
        }
        return props.Meta.display_name
    }

    const valueAsNumber = (val) => {
    }

    return (
        <Grid container alignItems="center" spacing={1}>
            <Grid item style={{minWidth: "100px", color: color()}}>
                <Tooltip title={tooltip()} TransitionComponent={Zoom}>
                    <div>{props.Meta.display_name}:</div>
                </Tooltip>
            </Grid>
            <Grid item>
                <Tooltip title={tooltip()} TransitionComponent={Zoom}>
                    <TextField
                        fullWidth={true}
                        value={parseFloat(props.Data)}
                        onChange={(ev) => {
                            let minValue = parseFloat(props.Meta.min_value)
                            let maxValue = parseFloat(props.Meta.max_value)
                            let valueAsNumber = parseFloat(ev.target.value)
                            if (valueAsNumber < minValue || valueAsNumber > maxValue) {
                                setOutOfRange(true)
                            } else {
                                setOutOfRange(false)
                            }
                            props.OnChangedValue(props.Meta.name, parseFloat(ev.target.value))
                        }}
                    />
                </Tooltip>
            </Grid>
        </Grid>
    )
}
