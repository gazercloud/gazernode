import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from '@material-ui/icons/Memory';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from "@material-ui/icons/Edit";
import Dialog from "@material-ui/core/Dialog";
import DialogTitle from "@material-ui/core/DialogTitle";
import DialogContent from "@material-ui/core/DialogContent";
import DialogContentText from "@material-ui/core/DialogContentText";
import DialogActions from "@material-ui/core/DialogActions";
import ArrowBackIosOutlinedIcon from "@material-ui/icons/ArrowBackIosOutlined";
import EditOutlinedIcon from '@material-ui/icons/EditOutlined';
import DeleteOutlinedIcon from '@material-ui/icons/DeleteOutlined';
import Divider from "@material-ui/core/Divider";
import WidgetError from "../Widgets/WidgetError";

function PageUnit(props) {
    const [unitValues, setUnitValues] = React.useState([])
    const [unitState, setUnitState] = React.useState([])
    const [unitName, setUnitName] = React.useState("")

    const [messageError, setMessageError] = React.useState("")
    const [stateLoading, setStateLoading] = React.useState(false)

    const btnStyle = (key) => {
        if (currentItem === key) {
            return {
                borderBottom: '1px solid #333333',
                cursor: "pointer",
                backgroundColor: "#222222",
            }
        } else {
            if (hoverItem === key) {
                return {
                    borderBottom: '1px solid #333333',
                    cursor: "pointer",
                    backgroundColor: "#222222"
                }
            } else {
                return {
                    borderBottom: '1px solid #333333',
                    cursor: "pointer",
                    backgroundColor: "#1E1E1E"
                }
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (id, name) => {
        const nameHex = new Buffer(name).toString('hex');
        props.OnNavigate("#form=data_item&dataItemName="+ nameHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestUnitItems = (unitName) => {
        let req = {
            name: unitName
        }
        Request('unit_items_values', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitValues(result)
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

    const requestRemoveUnit = (unitId) => {
        let req = {
            ids: [unitId]
        }
        Request('unit_remove', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            props.OnNavigate("#form=units")
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

    const requestUnitState = (unitId) => {
        if (stateLoading)
            return

        setStateLoading(true)

        let req = {
            id: unitId
        }
        Request('unit_state', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitState(result)
                            setUnitName(result.name)
                            props.OnTitleUpdate("Unit " + result.name)
                            requestUnitItems(result.name)
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
                    setStateLoading(false)
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

    const isServiceDataItemName = (dataItemName) => {
        if (dataItemName.includes("/.service/"))
            return true
        return false
    }

    const shortDataItemName = (dataItemName) => {
        dataItemName = dataItemName.replaceAll(unitName + "/", "")
        return dataItemName
    }

    const getUnitId = () => {
        return new Buffer(props.UnitId, 'hex').toString();
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        props.OnTitleUpdate("Unit ... loading")
        requestUnitState(getUnitId())
        setFirstRendering(false)
    }

    useEffect(() => {
        const unitId = new Buffer(props.UnitId, 'hex').toString();
        const timer = setInterval(() => {
            requestUnitState(unitId)
        }, 500);
        return () => clearInterval(timer);
    });

    const displayItemName = (item, mainItem) => {
        if (mainItem) {
            return (
                <div style={{fontSize: '24pt'}}>{shortDataItemName(item.name)}</div>)
        } else {
            return (
                <div style={{fontSize: '20pt'}}>{shortDataItemName(item.name)}</div>)
        }
    }

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
                        <span>{item.value.v} </span>
                        <span>{item.value.u}</span>
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
                        <span>{item.value.v} </span>
                        <span>{item.value.u}</span>
                    </div>
                    <div style={{fontSize: '10pt', color: '#AAA'}}>
                        <span>{dt2}</span>
                    </div>
                </div>
            )
        }
    }

    const displayItem = (item) => {
        if (item.name === unitState.main_item) {
            return (
                <div style={{color: '#C70'}}>
                    {displayItemName(item, true)}
                    {displayItemValue(item, true)}
                </div>
            )
        } else {
            return (
                <div>
                    {displayItemName(item)}
                    {displayItemValue(item)}
                </div>
            )
        }
    }

    let name = "unknown"
    let status = "unknown"

    if (unitState !== undefined) {
        name = unitState.name
        status = unitState.status
    }

    const displayItems = (onlyMain) => {
        if (unitValues === undefined || unitValues.items === undefined) {
            return (<div>loading ...</div>)
        } else {
            return (
                unitValues.items.map((item) => (
                    !isServiceDataItemName(item.name) && ((!onlyMain && unitState.main_item !== item.name) || (onlyMain && unitState.main_item === item.name)) ?
                        <ListItem
                            button
                            onMouseEnter={() => handleEnter(item.name)}
                            onMouseLeave={() => handleLeave(item.name)}
                            style={btnStyle(item.name)}
                            onClick={btnClick.bind(this, item.id, item.name)}
                        >
                            <ListItemIcon>
                                <MemoryIcon />
                            </ListItemIcon>
                            <div>
                                {displayItem(item)}
                            </div>
                        </ListItem>:<div/>
                ))
            )
        }
    }

    const btnClickUnitConfig = () => {
        const idHex = new Buffer(getUnitId()).toString('hex');
        props.OnNavigate("#form=unit_config&unitId=" + idHex)
    }

    const [dialogRemoveUnitVisible, setDialogRemoveUnitVisible] = React.useState(false)
    const dialogRemoveUnit = () => {
        return (
            <Dialog open={dialogRemoveUnitVisible} onClose={() => {
                setDialogRemoveUnitVisible(false)
            }
            } aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">Remove unit</DialogTitle>
                <DialogContent>
                    <DialogContentText>
                        Unit name: {unitName}
                    </DialogContentText>
                    <DialogContent>
                        Remove unit?
                    </DialogContent>
                </DialogContent>
                <DialogActions>
                    <Button onClick={()=>{ setDialogRemoveUnitVisible(false) }} color="primary">
                        Cancel
                    </Button>
                    <Button onClick={()=> {
                        requestRemoveUnit(getUnitId())
                    }} color="primary">
                        OK
                    </Button>
                </DialogActions>
            </Dialog>
        )
    }

    const btnClickUnitRemove = () => {
        setDialogRemoveUnitVisible(true)
    }

    return (
        <div>
            <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                <Tooltip title="Back" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={()=>{
                        window.history.back()
                    }} variant="outlined" color="primary">
                        <ArrowBackIosOutlinedIcon fontSize="large"  style={{borderBottom: "0px solid #00A0E3"}}/>
                    </IconButton>
                </Tooltip>

                <Divider orientation="vertical" flexItem style={{marginLeft: "20px", marginRight: "20px"}} />

                <Tooltip title="Edit" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={btnClickUnitConfig.bind(this)} variant="outlined" color="primary">
                        <EditOutlinedIcon fontSize="large"  style={{borderBottom: "0px solid #00A0E3"}}/>
                    </IconButton>
                </Tooltip>

                <Tooltip title="Remove unit" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                    <IconButton onClick={btnClickUnitRemove.bind(this)} variant="outlined" color="primary">
                        <DeleteOutlinedIcon fontSize="large"/>
                    </IconButton>
                </Tooltip>
            </Grid>

            <div style={{marginTop: '40px'}}>
                <Typography style={{fontSize: '10pt'}}>Unit</Typography>
                <Typography style={{fontSize: '24pt'}} color='primary'>{name}</Typography>
                <Typography style={{fontSize: '10pt'}}>Status</Typography>
                <Typography style={{fontSize: '16pt'}} color='primary'>{status}</Typography>
            </div>

            <Grid container direction="column">
                <Grid item>
                    <List component="nav">
                        {displayItems(true)}
                        {displayItems(false)}
                    </List>
                </Grid>
            </Grid>
            {dialogRemoveUnit()}
            <div>
                <WidgetError Message={messageError} />
            </div>
        </div>
    );
}

export default PageUnit;
