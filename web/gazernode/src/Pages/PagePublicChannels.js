import React, {useEffect, useState} from 'react';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Typography from "@material-ui/core/Typography";
import BlurOnIcon from "@material-ui/icons/BlurOn";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import PauseIcon from '@material-ui/icons/Pause';
import DialogAddPublicChannel from "../Dialogs/DialogAddPublicChannel";
import {useSnackbar} from "notistack";
import WidgetError from "../Widgets/WidgetError";

export default function PagePublicChannels(props) {
    const [channels, setChannels] = React.useState({})
    const { enqueueSnackbar } = useSnackbar();
    const [messageError, setMessageError] = React.useState("")

    const [stateLoading, setStateLoading] = React.useState(false)

    const btnStyle = (key) => {
        let borderTop = '0px solid #333333'
        let backColor = '#1E1E1E'
        let backColorHover = '#222222'


        if (currentItem === key) {
            return {
                marginRight: '10px',
                marginBottom: '10px',
                borderTop: borderTop,
                backgroundColor: backColor,
            }
        } else {
            if (hoverItem === key) {
                return {
                    padding: '20px',
                    minWidth: '300px',
                    maxWidth: '300px',
                    minHeight: '130px',
                    maxHeight: '130px',
                    borderRadius: '10px',
                    marginRight: '10px',
                    marginBottom: '10px',
                    borderTop: borderTop,
                    backgroundColor: backColor,
                }
            } else {
                return {
                    padding: '20px',
                    minWidth: '300px',
                    maxWidth: '300px',
                    minHeight: '130px',
                    maxHeight: '130px',
                    borderRadius: '10px',
                    marginRight: '10px',
                    marginBottom: '10px',
                    borderTop: borderTop,
                    backgroundColor: backColorHover,
                }
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (id, name) => {
        const idHex = new Buffer(id).toString('hex');
        props.OnNavigate("#form=public_channel&channelId="+ idHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestChannels = () => {
        if (stateLoading)
            return

        setStateLoading(true)

        let req = {
        }
        Request('public_channel_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setChannels(result)
                            setMessageError("")
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

    const requestStartChannel = (channelId) => {
        let req = {
            ids: [channelId]
        }
        Request('public_channel_start', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestChannels()
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

    const requestStopChannel = (channelId) => {
        let req = {
            ids: [channelId]
        }
        Request('public_channel_stop', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestChannels()
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
        props.OnTitleUpdate('Public Channels')
        requestChannels()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestChannels()
        }, 500);
        return () => clearInterval(timer);
    });

    const channelStatus = (item) => {
        if (item.status === 'started') {
            return (
                <div style={{color: "#080"}}>
                    started
                </div>
            );
        }

        if (item.status === 'stopped') {
            return (
                <div style={{color: "#777"}}>
                    stopped
                </div>
            )
        }

        return (
            <div style={{color: "#555"}}>-</div>
        )
    }

    const displayItem = (item) => {
        return (
            <Grid container direction='column'>
                <Grid item onClick={btnClick.bind(this, item.id, item.name)} style={{cursor: 'pointer'}}>
                    <Grid container>
                        <Grid item>
                            <Grid container alignItems='center'>
                                {item.status === 'started'?
                                    <BlurOnIcon color='primary' fontSize='large' style={{marginRight: '5px'}} /> :
                                    <BlurOnIcon color='disabled' fontSize='large' style={{marginRight: '5px'}} />
                                }

                                <Typography style={{
                                    fontSize: '14pt',
                                    overflow: 'hidden',
                                    height: '20pt',
                                    width: '200px',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                }}>{item.name}</Typography>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item style={{marginTop: '10px', backgroundColor: "#181818", borderRadius: "10px"}}>
                    <Grid container={true} direction="row" alignItems="center">
                        <Grid item style={{flexGrow: 1, paddingLeft: "10px"}}>
                            {channelStatus(item)}
                        </Grid>
                        <Grid item>
                            <Tooltip title="Start" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='secondary' style={{margin: '0px', opacity: "50%"}} disabled={item.status === 'started'}>
                                    <PlayArrowIcon fontSize="default" onClick={requestStartChannel.bind(this, item.id)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                        <Grid item>
                            <Tooltip title="Stop" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='primary' style={{margin: '0px', opacity: "50%"}} disabled={item.status === 'stopped'}>
                                    <PauseIcon fontSize="default" onClick={requestStopChannel.bind(this, item.id)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                    </Grid>

                </Grid>
            </Grid>
        )
    }

    return (
        <div>
            <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <DialogAddPublicChannel OnSuccess={()=> {
                    enqueueSnackbar("Channel created", {variant: "success"})
                    requestChannels()
                }} />
            </Grid>
            <Grid container direction="column" style={{marginTop: "20px"}}>
                <Grid item>
                    <Grid container direction="row">
                        {channels !== undefined && channels.channels !== undefined ? channels.channels.map((item) => (
                            <Grid item
                                  key={"gazer-page_channels-channel-" + item.id}
                                  button
                                  onMouseEnter={() => handleEnter(item.id)}
                                  onMouseLeave={() => handleLeave(item.id)}
                                  style={btnStyle(item.id)}
                            >
                                {displayItem(item)}
                            </Grid>
                        )) : <div/>}
                    </Grid>
                </Grid>
            </Grid>
            <div>
                <WidgetError Message={messageError} />
            </div>
        </div>
    );
}


