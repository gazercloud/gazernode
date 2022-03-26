import React, {useState} from 'react';
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import SearchIcon from "@material-ui/icons/Search";
import TextField from "@material-ui/core/TextField";
import Typography from "@material-ui/core/Typography";

export default function WidgetSensorTypeSelect(props) {

    const [sensorsTypesList, setSensorsTypesList] = React.useState([])
    const [filterString, setFilterString] = React.useState("")

    const onClickOnCard = (sType) => {
        props.onSelected(sType)
    }

    const requestSensorsTypes = (filterStr) => {
        let req = {
            "func": "host_get_sensors_types",
            "filter_string": filterStr
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            let sensorsTypes = []

                            for (let i = 0; i < result.sensors_types.length; i++)
                            {
                                let sItem = result.sensors_types[i]
                                sensorsTypes.push({
                                    type: sItem.type,
                                    name: sItem.name,
                                    desc: sItem.desc
                                })
                            }

                            setSensorsTypesList(sensorsTypes)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                //setProcessingLoadParameter(false)
            })
            .catch((err) => {
                //setProcessingLoadParameter(false)
                //setErrorMessage("Unknown error")
            })
    }


    const sensorCard = (type, name, desc) => {
        return (
            <Paper
                style={{padding: "10px", width: "300px", maxWidth: "300px", cursor: "pointer"}}
                onClick={onClickOnCard.bind(this, type)}
            >
                <Grid container direction="column">
                    <Grid item>
                        <Typography style={{fontSize: "16pt", color: "#333333"}}>
                            {type}
                        </Typography></Grid>
                    <Grid item>
                        <Typography style={{fontSize: "12pt", color: "#777777"}}>
                            {desc}
                        </Typography>
                    </Grid>
                </Grid>
            </Paper>
        )
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        requestSensorsTypes("")
        setFirstRendering(false)
    }


    return (
        <div>
            <Paper variant="outlined" style={{marginBottom: "10px", padding: "10px"}}>
                <Grid container spacing={1} alignItems="flex-end">
                    <Grid item>
                        <SearchIcon />
                    </Grid>
                    <Grid item>
                        <TextField
                            id="input-with-icon-grid"
                            label="Search sensor"
                            value={filterString}
                            onChange={(ev)=>{
                                setFilterString(ev.target.value)
                                requestSensorsTypes(ev.target.value)
                            }}
                            style={{width: "100%"}} />
                    </Grid>
                </Grid>
            </Paper>

            <Grid container direction="row" spacing={3}>
                {sensorsTypesList !== undefined ? sensorsTypesList.map((item) => (
                    <Grid item>
                        {sensorCard(item.type, item.name, item.desc)}
                    </Grid>
                )) : <div/>}
            </Grid>
        </div>
    );
}
