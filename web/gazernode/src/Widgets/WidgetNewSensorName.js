import React, {useState} from 'react';
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import Button from "@material-ui/core/Button";

export default function WidgetNewSensorName(props) {
    const [sensorName, setSensorName] = useState("")

    return (
        <div>
            <TextField label="Sensor name" value={sensorName}
                       onChange={(ev)=>{setSensorName(ev.target.value)}}
            />
            <Grid container direction="row">
                <Grid item>
                    <Button
                        onClick={() => {props.Cancel()}}>
                        Back
                    </Button>
                </Grid>
                <Grid item>
                    <Button
                        onClick={() => {
                            props.OK(sensorName)
                        }}>
                        Next
                    </Button>
                </Grid>
            </Grid>

        </div>
    );
}
