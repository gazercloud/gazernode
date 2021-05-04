import React, {useState} from 'react';
import WidgetDataItems from "../Widgets/WidgetDataItems";
import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import WidgetDataItemDetails from "../Widgets/WidgetDataItemDetails";

export default function PageDataItems(props) {
    const [currentDataItemPath, setCurrentDataItemPath] = useState("")

    return (
        <Grid container direction="row">
            <Grid item>
                <Paper variant="outlined" style={{maxWidth: "700px"}}>
                    <WidgetDataItems
                        Root=""
                        OnDataItemSelected={(path) => {
                            setCurrentDataItemPath(path)
                        }}
                    />
                </Paper>
            </Grid>
            <Grid item>
                <Paper variant="outlined" style={{marginLeft: "5px", padding: "5px",minWidth: "400px", maxWidth: "400px", minHeight: "400px"}} >
                    <WidgetDataItemDetails Path={currentDataItemPath}  ShowSensorInfo={true}/>
                </Paper>
            </Grid>
        </Grid>
    );
}
