import React, {useState} from 'react';
import TreeItem from '@material-ui/lab/TreeItem';
import TreeView from '@material-ui/lab/TreeView';
import SupervisorAccountIcon from '@material-ui/icons/SupervisorAccount';
import InfoIcon from '@material-ui/icons/Info';
import ForumIcon from '@material-ui/icons/Forum';
import LocalOfferIcon from '@material-ui/icons/LocalOffer';
import ArrowDropDownIcon from '@material-ui/icons/ArrowDropDown';
import ArrowRightIcon from '@material-ui/icons/ArrowRight';
import MailIcon from '@material-ui/icons/Mail';
import DeleteIcon from '@material-ui/icons/Delete';
import Typography from "@material-ui/core/Typography";
import Label from '@material-ui/icons/Label';
import { makeStyles } from '@material-ui/core/styles';
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
