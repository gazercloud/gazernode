import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import BlurOnIcon from "@material-ui/icons/BlurOn";
import {Fab, Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from '@material-ui/icons/Edit';
import PlayArrowIcon from '@material-ui/icons/PlayArrow';
import PauseIcon from '@material-ui/icons/Pause';
import AddIcon from '@material-ui/icons/Add';
import AddOutlinedIcon from "@material-ui/icons/AddOutlined";
import WidgetError from "../Widgets/WidgetError";

function PageUnits(props) {
    const [units, setUnits] = React.useState({})
    const [messageError, setMessageError] = React.useState("")
    const [stateLoading, setStateLoading] = React.useState(false)
    const [hoverItem, setHoverItem] = useState("")

    const btnStyle = (key) => {
        let borderTop = '0px solid #333333'
        let backColor = '#1E1E1E'
        let backColorHover = '#222222'


        if (hoverItem === key) {
            return {
                padding: '20px',
                minWidth: '300px',
                maxWidth: '300px',
                minHeight: '200px',
                maxHeight: '200px',
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
                minHeight: '200px',
                maxHeight: '200px',
                borderRadius: '10px',
                marginRight: '10px',
                marginBottom: '10px',
                borderTop: borderTop,
                backgroundColor: backColorHover,
            }
        }
    }


    const btnClick = (id, name) => {
        const idHex = new Buffer(id).toString('hex');
        props.OnNavigate("#form=unit&unitId="+ idHex)
    }

    const btnClickUnitConfig = (id) => {
        const idHex = new Buffer(id).toString('hex');
        props.OnNavigate("#form=unit_config&unitId=" + idHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestUnitsState = () => {
        if (stateLoading)
            return

        setStateLoading(true)

        let req = {
        }
        Request('unit_state_all', req)
            .then((res) => {

                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnits(result)
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
            setStateLoading(false)
            RequestFailed()
        });
    }

    const requestStartUnit = (unitId) => {
        let req = {
            ids: [unitId]
        }
        Request('unit_start', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUnitsState()
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                    RequestFailed()
                }
            });
    }

    const requestStopUnit = (unitId) => {
        let req = {
            ids: [unitId]
        }
        Request('unit_stop', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUnitsState()
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                    RequestFailed()
                }
            });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate('Units')
        requestUnitsState()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestUnitsState()
        }, 500);
        return () => clearInterval(timer);
    });

    const displayUnitValue = (unit) => {
        if (unit.status !== "started") {
            return (
                <div style={{
                    fontSize: '24pt',
                    color: '#555',
                    overflow: 'hidden',
                    height: '30pt',
                    width: '260px',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    textAlign: 'center'
                }}>
                    {unit.value + "  " + unit.uom}
                </div>
            )
        }

        if (unit.uom === "error") {
            return (
                <div style={{
                    fontSize: '24pt',
                    color: '#F30',
                    overflow: 'hidden',
                    height: '30pt',
                    width: '260px',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    textAlign: 'center'
                }}>
                    {unit.value + " " + unit.uom}
                </div>
            )
        }

        return (
            <div style={{
                fontSize: '24pt',
                color: '#080',
                overflow: 'hidden',
                height: '30pt',
                width: '260px',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                textAlign: 'center'
            }}>
                {unit.value + "  " + unit.uom}
            </div>
        )
    }

    const unitStatus = (item) => {
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
                        <Grid item style={{marginTop: '10px'}}>
                                {displayUnitValue(item)}
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item style={{marginTop: '35px', backgroundColor: "#181818", borderRadius: "10px"}}>
                    <Grid container={true} direction="row" alignItems="center">
                        <Grid item>
                            <Tooltip title="Edit unit" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='primary' style={{margin: '0px', opacity: "50%"}}>
                                    <EditIcon fontSize="default" onClick={btnClickUnitConfig.bind(this, item.id, item.name)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                        <Grid item style={{flexGrow: 1}}>
                            {unitStatus(item)}
                        </Grid>
                        <Grid item>
                            <Tooltip title="Start" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='secondary' style={{margin: '0px', opacity: "50%"}} disabled={item.status === 'started'}>
                                    <PlayArrowIcon fontSize="default" onClick={requestStartUnit.bind(this, item.id)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                        <Grid item>
                            <Tooltip title="Stop" TransitionComponent={Zoom}>
                                <IconButton variant='outlined' color='primary' style={{margin: '0px', opacity: "50%"}} disabled={item.status === 'stopped'}>
                                    <PauseIcon fontSize="default" onClick={requestStopUnit.bind(this, item.id)} />
                                </IconButton>
                            </Tooltip>
                        </Grid>
                    </Grid>

                </Grid>
            </Grid>
        )
    }

    const onClickAddUnit = () => {
        props.OnNavigate("#form=unit_add")
    }

    return (
        <div>
            <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <Tooltip title="Add unit" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={onClickAddUnit} variant="outlined" color="primary">
                        <AddOutlinedIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>
            </Grid>
            <Grid container direction="column" style={{marginTop: "20px"}}>
                <Grid item>
                    <Grid container direction="row">
                            {units !== undefined && units.items !== undefined ? units.items.map((item) => (
                                <Grid item
                                    key={"gazer-page_units-unit-" + item.id}
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

export default PageUnits;
