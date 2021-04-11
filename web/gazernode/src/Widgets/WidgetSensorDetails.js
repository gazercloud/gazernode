import React, {useState} from 'react';
import Paper from "@material-ui/core/Paper";
import WidgetDataItems from "./WidgetDataItems";
import WidgetDataItemDetails from "./WidgetDataItemDetails";
import Grid from "@material-ui/core/Grid";
import WidgetSensorInfo from "./WidgetSensorInfo";

function SensorsDetails(props) {
    const [sensorId, setSensorID] = useState("")
    const [sensorName, setSensorName] = useState("")
    const [currentDataItemPath, setCurrentDataItemPath] = useState("")


    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering || sensorId !== props.SensorID) {
        setSensorID(props.SensorID)
        setSensorName(props.SensorName)
        setFirstRendering(false)
        setCurrentDataItemPath("")
    }

    return (
        <Paper variant="outlined" style={{padding: "5px"}}>
            <Grid container direction="row">
                <Grid item>
                    <Paper variant="outlined"
                           style={{padding: "5px", marginBottom: "10px", minWidth: "600px", maxWidth: "600px"}}>
                        <WidgetSensorInfo SensorID={sensorId} SensorName={sensorName}/>
                    </Paper>
                    <Paper variant="outlined" style={{padding: "5px", minWidth: "600px", maxWidth: "600px"}}>
                        <WidgetDataItems
                            Root={props.SensorName}
                            OnDataItemSelected={(path) => {
                                setCurrentDataItemPath(path)
                            }}
                        />
                    </Paper>
                </Grid>
                <Grid item>
                    <Paper variant="outlined" style={{
                        marginLeft: "5px",
                        padding: "5px",
                        minWidth: "400px",
                        maxWidth: "400px",
                        minHeight: "400px"
                    }}>
                        <WidgetDataItemDetails Path={currentDataItemPath}/>
                    </Paper>
                </Grid>
            </Grid>
        </Paper>
    );

}

export default SensorsDetails;
