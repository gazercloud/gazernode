import React, {useState} from 'react';
import {makeStyles} from "@material-ui/core/styles";
import Grid from "@material-ui/core/Grid";
import Paper from "@material-ui/core/Paper";
import SearchIcon from "@material-ui/icons/Search";
import TextField from "@material-ui/core/TextField";
import Button from "@material-ui/core/Button";
import AddBoxOutlinedIcon from "@material-ui/icons/AddBoxOutlined";
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from "@material-ui/icons/Memory";
import ListItemText from "@material-ui/core/ListItemText";
import Tabs from "@material-ui/core/Tabs";
import Tab from "@material-ui/core/Tab";
import WidgetUserProperties from "../Widgets/WidgetUserProperties";

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
}));

function PageUsers(props) {
    const classes = useStyles();
    const [currentTabIndex, setCurrentTabIndex] = React.useState(0)
    const [addSensorInterfaceVisible, setAddSensorInterfaceVisible] = React.useState(false)
    const [usersList, setUsersList] = React.useState([])
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
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (ev, key) => {
        setCurrentItem(ev)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }


    const requestUsers = (filterStr) => {
        let req = {
            "fn": "host_get_users",
            "filter_string": filterStr
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            let users = []

                            for (let i = 0; i < result.users.length; i++)
                            {
                                let sItem = result.users[i]
                                users.push({
                                    id: sItem.id,
                                    name: sItem.name
                                })
                            }

                            setUsersList(users)
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

    const requestAddUser = () => {
        let req = {
            "fn": "host_add_user",
            "name": "user1",
            "password": "password1"
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUsers(filterString)
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

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        requestUsers(filterString)
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
                                    label="Search user ..."
                                    onChange={(ev)=>{
                                        setCurrentItem("")
                                        setFilterString(ev.target.value)
                                        requestUsers(ev.target.value)
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
                            requestAddUser()
                        }}
                    >
                        ADD USER
                    </Button>
                    <Paper variant="outlined" className={classes.root} style={{marginTop: "10px"}}>
                        <List component="nav" aria-label="main mailbox folders">
                            {usersList !== undefined ? usersList.map((item) => (
                                <ListItem
                                    button
                                    onMouseEnter={() => handleEnter(item.id)}
                                    onMouseLeave={() => handleLeave(item.id)}
                                    style={btnStyle(item.id)}
                                    onClick={btnClick.bind(this, item.id)}
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
                        <Tab label="Properties"/>
                        <Tab label="Group's member"/>
                    </Tabs>
                    <div hidden={currentTabIndex !== 0}>
                        <WidgetUserProperties CurrentId={currentItem} OnNeedToLoadUserList={() => {
                            requestUsers(filterString)
                            setCurrentItem("")
                        }} />
                    </div>
                    <div hidden={currentTabIndex !== 1}>
                        Groups
                    </div>
                </Grid>
            </Grid>
        </div>
    );
}

export default PageUsers;
