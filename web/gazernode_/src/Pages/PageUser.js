import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from '@material-ui/icons/Memory';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import {Paper, Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from "@material-ui/icons/Edit";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogActions from "@material-ui/core/DialogActions";
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined';
import ArrowBackIosOutlinedIcon from "@material-ui/icons/ArrowBackIosOutlined";
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import PauseIcon from '@material-ui/icons/Pause';
import DialogEditPublicChannel from "../Dialogs/DialogEditPublicChannel";
import Divider from "@material-ui/core/Divider";
import HelpOutlineOutlinedIcon from '@material-ui/icons/HelpOutlineOutlined';
import DialogEditUser from "../Dialogs/DialogEditUser";
import WidgetError from "../Widgets/WidgetError";

export default function PageUser(props) {
    const [sessions, setSessions] = React.useState({})
    const [sessionsLoaded, setSessionsLoaded] = React.useState(false)
    const [messageError, setMessageError] = React.useState("")

    const btnStyle = (key) => {
        if (hoverItem === key) {
            return {
                borderBottom: '1px solid #333333',
                backgroundColor: "#222222"
            }
        } else {
            return {
                borderBottom: '1px solid #333333',
                backgroundColor: "#1E1E1E"
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (id, name) => {
        const nameHex = new Buffer(name).toString('hex');
        props.OnNavigate("#form=data_item&dataItemName=" + nameHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestRemoveUser = () => {
        let req = {
            user_name: getUserName()
        }
        Request('user_remove', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            props.OnNavigate("#form=users")
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

    const requestRemoveSession = (sessionToken) => {
        let req = {
            session_token: sessionToken
        }
        Request('session_remove', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUserSessions(getUserName())
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

    const requestUserSessions = () => {
        let req = {
            user_name: getUserName()
        }
        Request('session_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            if (result.items === null || result.items === undefined) {
                                result.items = []
                            }
                            setSessions(result)
                            setSessionsLoaded(true)
                            setMessageError("")
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                        }
                    );
                    RequestFailed()
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
        });
    }


    const isServiceDataItemName = (dataItemName) => {
        if (dataItemName.includes("/.service/"))
            return true
        return false
    }

    const getUserName = () => {
        return new Buffer(props.UserName, 'hex').toString();
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("User " + getUserName())
        requestUserSessions(getUserName())
        setFirstRendering(false)
    }

    useEffect(() => {
        const userName = getUserName()

        const timer = setInterval(() => {
            requestUserSessions(userName)
        }, 2000);
        return () => clearInterval(timer);
    });

    const displayItemValue = (item, mainItem) => {
        let fontSize = '14pt'
        if (mainItem) {
            fontSize = '24pt'
        }

        return (
            <div>
                <div style={{fontSize: fontSize, color: '#080'}}>
                    <span>{item.session_token} </span>
                </div>
                <div style={{fontSize: '10pt', color: '#AAA'}}>
                </div>
            </div>
        )
    }

    const displayItem = (item) => {
        let dt = new Date(item.session_open_time / 1000)
        let dd = String(dt.getDate()).padStart(2, '0');
        let mm = String(dt.getMonth() + 1).padStart(2, '0'); //January is 0!
        let yyyy = dt.getFullYear();
        let hh = dt.getHours();
        if (hh < 10) hh = '0' + hh;

        let min = dt.getMinutes();
        if (min < 10) min = '0' + min;

        let ss = dt.getSeconds();
        if (ss < 10) ss = '0' + ss;

        let today = yyyy + '-' + mm + '-' + dd + " " + hh + ":" + min + ":" + ss;

        return (
            <div>
                <div style={{fontSize: '14pt'}}>{today}</div>
                <div style={{fontSize: '10pt', color: "#444"}}>{item.session_token}</div>
            </div>
        )
    }

    const displayItems = () => {
        if (!sessionsLoaded) {
            return (<div>loading ...</div>)
        } else {
            if (sessions.items.length < 1) {
                return (<div>no sessions</div>)
            } else {
                return (
                    sessions.items.map((item) => (
                        <ListItem
                            onMouseEnter={() => handleEnter(item.session_token)}
                            onMouseLeave={() => handleLeave(item.session_token)}
                            style={btnStyle(item.session_token)}
                        >
                            <ListItemIcon>
                                <IconButton
                                    onClick={btnClickItemRemove.bind(this, item.session_token)}><DeleteOutlinedIcon/></IconButton>
                            </ListItemIcon>
                            <div>
                                {displayItem(item)}
                            </div>
                        </ListItem>
                    ))
                )
            }
        }
    }

    const [dialogRemoveUserVisible, setDialogRemoveUserVisible] = React.useState(false)
    const dialogRemoveChannel = () => {
        return (
            <Dialog open={dialogRemoveUserVisible} onClose={() => {
                setDialogRemoveUserVisible(false)
            }
            } aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">
                    <Grid container alignItems="center">
                        <HelpOutlineOutlinedIcon fontSize="large" style={{marginRight: "10px"}}/>
                        <span  style={{margin: "0px"}}>Confirmation</span>
                    </Grid>
                </DialogTitle>
                <DialogContent style={{backgroundColor: "#222", padding: "20px"}}>
                    <Typography style={{fontSize: "16pt", minWidth: "300px"}}>User Name: <span style={{color: "#F40"}}>{getUserName()}</span></Typography>
                    <Typography style={{fontSize: "16pt", minWidth: "300px", color: "#F40"}}>Remove user?</Typography>
                </DialogContent>
                <DialogActions style={{backgroundColor: "#222"}}>
                    <Button onClick={() => {
                        requestRemoveUser(getUserName())
                    }} color="primary"  variant="contained" style={{minWidth: "100px"}}>
                        OK
                    </Button>
                    <Button onClick={() => {
                        setDialogRemoveUserVisible(false)
                    }} color="primary" variant="outlined" style={{minWidth: "100px"}}>
                        Cancel
                    </Button>
                </DialogActions>
            </Dialog>
        )
    }

    const btnClickUserRemove = () => {
        setDialogRemoveUserVisible(true)
    }
    const btnClickItemRemove = (token) => {
        requestRemoveSession(token)
    }

    return (
        <div>

            <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <Tooltip title="Back" TransitionComponent={Zoom} style={{border: "1px solid #333"}}>
                    <IconButton onClick={() => {
                        window.history.back()
                    }} variant="outlined" color="primary">
                        <ArrowBackIosOutlinedIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>

                <Divider orientation="vertical" flexItem style={{marginLeft: "20px", marginRight: "20px"}} />

                <DialogEditUser UserName={getUserName()}/>

                <Tooltip title="Remove user" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={btnClickUserRemove.bind(this)} variant="outlined" color="primary">
                        <DeleteOutlinedIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>
            </Grid>

            <div style={{marginTop: '40px'}}>
                <Typography style={{fontSize: '10pt'}}>User Name:</Typography>
                <Typography style={{fontSize: '24pt'}} color='primary'>{getUserName()}</Typography>
            </div>

            <Typography style={{fontSize: '10pt'}}>Sessions:</Typography>
            <Grid container direction="column" >
                <Grid item>
                    <List component="nav">
                        {displayItems(false)}
                    </List>
                </Grid>
            </Grid>
            {dialogRemoveChannel()}
            <div>
                <WidgetError Message={messageError} />
            </div>
        </div>
    );
}

