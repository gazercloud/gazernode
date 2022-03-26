import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogActions from "@material-ui/core/DialogActions";
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined';
import ArrowBackIosOutlinedIcon from "@material-ui/icons/ArrowBackIosOutlined";
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import PauseIcon from '@material-ui/icons/Pause';
import DialogEditPublicChannel from "../Dialogs/DialogEditPublicChannel";
import Divider from "@material-ui/core/Divider";
import HelpOutlineOutlinedIcon from '@material-ui/icons/HelpOutlineOutlined';
import WidgetError from "../Widgets/WidgetError";

export default function PagePublicChannel(props) {
    const [channelValues, setChannelValues] = React.useState([])
    const [channelState, setChannelState] = React.useState([])
    const [channelName, setChannelName] = React.useState("")

    const [btnStartDisable, setBtnStartDisable] = React.useState(true)
    const [btnStopDisable, setBtnStopDisable] = React.useState(true)

    const [messageError, setMessageError] = React.useState("")

    const [stateLoading, setStateLoading] = React.useState(false)

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

    const requestChannelItems = (channelId) => {
        let req = {
            id: channelId
        }
        Request('public_channel_item_state', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setChannelValues(result)
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
                    console.log("FAIL 1")
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
            console.log("FAIL 2")
        });
    }

    const requestRemoveChannel = (channelId) => {
        let req = {
            id: channelId
        }
        Request('public_channel_remove', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            props.OnNavigate("#form=public_channels")
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

    const requestRemoveItem = (channelId, itemName) => {
        let req = {
            ids: [channelId],
            items: [itemName]
        }
        Request('public_channel_item_remove', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestChannelItems(getChannelId())
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

    const requestChannelState = () => {
        if (stateLoading)
            return

        setStateLoading(true)

        let req = {}
        Request('public_channel_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            for (let i = 0; i < result.channels.length; i++) {
                                let ch = result.channels[i]
                                if (ch.id === getChannelId()) {
                                    setChannelState(ch)
                                    setChannelName(ch.name)

                                    if (ch.status === "started") {
                                        setBtnStartDisable(true)
                                    } else {
                                        setBtnStartDisable(false)
                                    }

                                    if (ch.status === "stopped") {
                                        setBtnStopDisable(true)
                                    } else {
                                        setBtnStopDisable(false)
                                    }

                                    props.OnTitleUpdate("Public Channel " + ch.name)
                                    requestChannelItems(getChannelId())
                                }
                            }

                            setStateLoading(false)
                        }
                    );
                    return
                }

                if (res.status === 500) {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                            setStateLoading(false)
                        }
                    );
                    RequestFailed()
                    return
                }

                setMessageError("Error: " + res.status + " " + res.statusText)
                RequestFailed()
                setStateLoading(false)

            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
            setStateLoading(false)
        });
    }

    const requestStartChannel = () => {
        let req = {
            ids: [getChannelId()]
        }
        Request('public_channel_start', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestChannelState()
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

    const requestStopChannel = () => {
        let req = {
            ids: [getChannelId()]
        }
        Request('public_channel_stop', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestChannelState()
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

    const getChannelId = () => {
        return new Buffer(props.ChannelId, 'hex').toString();
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("Public Channel ... loading")
        requestChannelState(getChannelId())
        setFirstRendering(false)
    }

    useEffect(() => {
        const channelId = new Buffer(props.ChannelId, 'hex').toString();
        const timer = setInterval(() => {
            requestChannelState(channelId)
        }, 1000);
        return () => clearInterval(timer);
    });

    const formatTime = (date) => {

        let hh = date.getHours();
        if (hh < 10) hh = '0' + hh;

        let mm = date.getMinutes();
        if (mm < 10) mm = '0' + mm;

        let ss = date.getSeconds();
        if (ss < 10) ss = '0' + ss;

        let fff = date.getMilliseconds();
        if (fff < 10) {
            fff = '00' + fff;
        } else {
            if (fff < 100)
                fff = '0' + fff
        }

        return hh + ':' + mm + ':' + ss + '.' + fff;
    }

    const displayItemValue = (item, mainItem) => {
        let dt1 = new Date(item.value.t / 1000)
        let dt2 = ""
        if (dt1 !== undefined) {
            dt2 = formatTime(dt1)
        }

        let fontSize = '14pt'
        if (mainItem) {
            fontSize = '24pt'
        }

        if (item.value.u === "error") {
            return (
                <div>
                    <div style={{fontSize: fontSize, color: '#F30'}}>
                        <span style={{color: stoppedColor()}}>{item.value.v} </span>
                        <span style={{color: stoppedColor()}}>{item.value.u}</span>
                    </div>
                    <div style={{fontSize: '10pt', color: '#AAA'}}>
                        <span>{dt2}</span>
                    </div>
                </div>
            )
        } else {
            return (
                <div>
                    <div style={{fontSize: fontSize, color: '#080'}}>
                        <span style={{color: stoppedColor()}}>{item.value.v} </span>
                        <span style={{color: stoppedColor()}}>{item.value.u}</span>
                    </div>
                    <div style={{fontSize: '10pt', color: '#AAA'}}>
                        <span style={{color: stoppedColor()}}>{dt2}</span>
                    </div>
                </div>
            )
        }
    }

    const displayItem = (item) => {
        return (
            <div onClick={btnClick.bind(this, item.id, item.name)} style={{cursor: "pointer"}}>
                <div style={{fontSize: '14pt'}}>{item.name}</div>
                {displayItemValue(item)}
            </div>
        )
    }

    let name = "unknown"
    let status = "unknown"

    if (channelState !== undefined) {
        name = channelState.name
        status = channelState.status
    }

    const displayItems = () => {
        if (channelValues === undefined || channelValues.items === undefined) {
            return (<div>loading ...</div>)
        } else {
            return (
                channelValues.items.map((item) => (
                    <ListItem
                        onMouseEnter={() => handleEnter(item.name)}
                        onMouseLeave={() => handleLeave(item.name)}
                        style={btnStyle(item.name)}
                    >
                        <ListItemIcon>
                            <IconButton
                                onClick={btnClickItemRemove.bind(this, item.name)}><DeleteOutlinedIcon/></IconButton>
                        </ListItemIcon>
                        <div>
                            {displayItem(item)}
                        </div>
                    </ListItem>
                ))
            )
        }
    }

    const [dialogRemoveChannelVisible, setDialogRemoveChannelVisible] = React.useState(false)
    const dialogRemoveChannel = () => {
        return (
            <Dialog open={dialogRemoveChannelVisible} onClose={() => {
                setDialogRemoveChannelVisible(false)
            }
            } aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">
                    <Grid container alignItems="center">
                        <HelpOutlineOutlinedIcon fontSize="large" style={{marginRight: "10px"}}/>
                        <span  style={{margin: "0px"}}>Confirmation</span>
                    </Grid>
                </DialogTitle>
                <DialogContent style={{backgroundColor: "#222", padding: "20px"}}>
                    <Typography style={{fontSize: "16pt", minWidth: "300px"}}>Channel Name: <span style={{color: "#F40"}}>{channelName}</span></Typography>
                    <Typography style={{fontSize: "16pt", minWidth: "300px", color: "#F40"}}>Remove channel?</Typography>
                </DialogContent>
                <DialogActions style={{backgroundColor: "#222"}}>
                    <Button onClick={() => {
                        requestRemoveChannel(getChannelId())
                    }} color="primary"  variant="contained" style={{minWidth: "100px"}}>
                        OK
                    </Button>
                    <Button onClick={() => {
                        setDialogRemoveChannelVisible(false)
                    }} color="primary" variant="outlined" style={{minWidth: "100px"}}>
                        Cancel
                    </Button>
                </DialogActions>
            </Dialog>
        )
    }

    const btnClickChannelRemove = () => {
        setDialogRemoveChannelVisible(true)
    }
    const btnClickItemRemove = (itemName) => {
        requestRemoveItem(getChannelId(), itemName)
    }

    const channelUrl = () => {
        return "https://gazer.cloud/channel/" + getChannelId()
    }

    const stoppedColor = () => {
        if (channelState.status === "started")
            return null
        else
            return "#444"
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

                <Tooltip title="Start channel" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton disabled={btnStartDisable} onClick={requestStartChannel} variant="outlined" color="primary">
                        <PlayArrowIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Stop channel" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton disabled={btnStopDisable} onClick={requestStopChannel} variant="outlined"
                                color="primary">
                        <PauseIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>

                <Divider orientation="vertical" flexItem style={{marginLeft: "20px", marginRight: "20px"}} />

                <DialogEditPublicChannel ChannelId={getChannelId()}/>

                <Tooltip title="Remove channel" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={btnClickChannelRemove.bind(this)} variant="outlined" color="primary">
                        <DeleteOutlinedIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>
            </Grid>

            <div style={{marginTop: '40px'}}>
                <Typography style={{fontSize: '10pt'}}>Channel ID</Typography>
                <Typography style={{fontSize: '24pt'}} color='primary'>
                    <a style={{color: "#00A0E3"}} target="_blank" href={channelUrl()}>{getChannelId()}</a>
                </Typography>
                <Typography style={{fontSize: '10pt'}}>Channel Name</Typography>
                <Typography style={{fontSize: '24pt', color: stoppedColor()}} color='primary'>{name}</Typography>
                <Typography style={{fontSize: '10pt'}}>Status</Typography>
                <Typography style={{fontSize: '16pt', color: stoppedColor()}} color='primary'>{status}</Typography>
            </div>

            <Grid container direction="column" >
                <Grid item>
                    <List component="nav" style={{color: stoppedColor()}}>
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

