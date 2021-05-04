import React, {useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemText from "@material-ui/core/ListItemText";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import Paper from "@material-ui/core/Paper";
import MemoryIcon from '@material-ui/icons/Memory';
import SearchIcon from '@material-ui/icons/Search';
import TextField from "@material-ui/core/TextField";
import Grid from "@material-ui/core/Grid";
import SensorsDetails from "../Widgets/WidgetSensorDetails";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import WidgetSensorsConfiguration from "../Widgets/WidgetSensorConfiguration";
import Button from "@material-ui/core/Button";
import AddBoxOutlinedIcon from '@material-ui/icons/AddBoxOutlined';
import Request from "../request";

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
}));

function PageSensorsOld(props) {
    const classes = useStyles();
    const [currentTabIndex, setCurrentTabIndex] = React.useState(0)
    const [addSensorInterfaceVisible, setAddSensorInterfaceVisible] = React.useState(false)
    const [sensorsList, setSensorsList] = React.useState([])
    const [filterString, setFilterString] = React.useState("")
    const [currentDataItemPath, setCurrentDataItemPath] = React.useState("")

    const tabChanged = (event, newValue) => {
        setCurrentTabIndex(newValue)
    }

    const btnStyle = (key) => {
        if (currentItem === key) {
            return {
                cursor: "pointer",
                backgroundColor: "#52bdff",
            }
        } else {
            if (hoverItem === key) {
                return {
                    cursor: "pointer",
                    backgroundColor: "#b8e4ff"
                }
            } else {
                return {
                    cursor: "pointer",
                    backgroundColor: "#FFFFFF"
                }
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [currentItemName, setCurrentItemName] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (ev, name) => {
        setCurrentItem(ev)
        setCurrentItemName(name)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }


    const requestSensors = (filterStr) => {
        let req = {
        }
        Request('unit_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setSensorsList(result)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        requestSensors(filterString)
        setFirstRendering(false)
    }

    return (
        <div>
            <Grid container direction="row" alignItems="flex-start" spacing={3}>
                <Grid item style={{minWidth: "400px", maxWidth: "400px"}}>
                    <Paper variant="outlined" className={classes.root} style={{marginBottom: "10px", padding: "10px"}}>
                        <Grid container spacing={1} alignItems="flex-end">
                            <Grid item>
                                <SearchIcon />
                            </Grid>
                            <Grid item>
                                <TextField
                                    id="input-with-icon-grid"
                                    value={filterString}
                                    label="Search sensor"
                                    onChange={(ev)=>{
                                        setCurrentItem("")
                                        setFilterString(ev.target.value)
                                        requestSensors(ev.target.value)
                                    }}

                                    style={{minWidth: "300px"}}
                                />
                            </Grid>
                        </Grid>
                    </Paper>
                    <Button
                        variant="contained"
                        startIcon={<AddBoxOutlinedIcon />}
                        onClick={() => {
                            props.onAddSensor()
                        }}
                    >
                        ADD SENSOR
                    </Button>
                    <Paper variant="outlined" className={classes.root} style={{marginTop: "10px"}}>
                        <List component="nav" aria-label="main mailbox folders">
                            {sensorsList !== undefined && sensorsList.items !== undefined ? sensorsList.items.map((item) => (
                                <ListItem
                                    button
                                    onMouseEnter={() => handleEnter(item.id)}
                                    onMouseLeave={() => handleLeave(item.id)}
                                    style={btnStyle(item.id)}
                                    onClick={btnClick.bind(this, item.id, item.name)}
                                >
                                    <ListItemIcon>
                                        <MemoryIcon />
                                    </ListItemIcon>
                                    <ListItemText primary={item.name} secondary={item.id} />
                                </ListItem>
                            )) : <div/>}
                        </List>
                    </Paper>
                </Grid>
                <Grid item>
                    <Tabs value={currentTabIndex} onChange={tabChanged} indicatorColor="primary" textColor="primary">
                        <Tab label="Data Items"/>
                        <Tab label="Configuration"/>
                   </Tabs>
                    <div hidden={currentTabIndex !== 0}>
                        <SensorsDetails SensorID={currentItem} SensorName={currentItemName}/>
                    </div>
                    <div hidden={currentTabIndex !== 1}>
                        <WidgetSensorsConfiguration SensorID={currentItem}/>
                    </div>
                </Grid>
            </Grid>
        </div>
    );
}

export default PageSensorsOld;
