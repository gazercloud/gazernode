import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from '@material-ui/icons/Memory';
import Grid from "@material-ui/core/Grid";
import Request, {RequestFailed} from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import ConfigSensor from "../Widgets/SensorConfig/ConfigSensor";
import TextField from "@material-ui/core/TextField";
import BlurOnIcon from "@material-ui/icons/BlurOn";
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import EditIcon from "@material-ui/icons/Edit";
import PlayArrowIcon from "@material-ui/icons/PlayArrow";
import PauseIcon from "@material-ui/icons/Pause";
import ArrowBackIosOutlinedIcon from "@material-ui/icons/ArrowBackIosOutlined";

function PageAddUnit(props) {
    const [filter, setFilter] = React.useState("")
    const [unitTypes, setUnitTypes] = React.useState("")

    const [messageError, setMessageError] = React.useState("")

    const btnStyle = (key) => {
        let borderTop = '0px solid #333333'
        let backColor = '#1E1E1E'
        let backColorHover = '#222222'


        if ("123" === key) {
            return {
                margin: '10px',
                borderTop: borderTop,
                backgroundColor: backColor,
            }
        } else {
            if (hoverItem === key) {
                return {
                    padding: '20px',
                    minWidth: '300px',
                    maxWidth: '300px',
                    minHeight: '200px',
                    maxHeight: '200px',
                    borderRadius: '10px',
                    margin: '10px',
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
                    margin: '10px',
                    borderTop: borderTop,
                    backgroundColor: backColorHover,
                }
            }
        }
    }

    const [hoverItem, setHoverItem] = useState("")

    const requestSensorTypes = () => {
        let req = {
            category: "",
            filter: "",
            offset: 0,
            max_count: 100
        }
        Request('unit_type_list', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitTypes(result)
                            setMessageError("")
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
        });
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        requestSensorTypes()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
        }, 500);
        return () => clearInterval(timer);
    });

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const btnClick = (type, name) => {
        const typeHex = new Buffer(type).toString('hex');
        props.OnNavigate("#form=unit_config&unitType="+ typeHex)
    }

    const displayUnitValue = (unit) => {
        return (
            <div style={{
                fontSize: '12pt',
                color: '#555',
                overflow: 'hidden',
                height: '50pt',
                width: '260px',
                textOverflow: 'ellipsis',
                textAlign: 'center'
            }}>
                {unit.help}
            </div>
        )
    }

    const displayItem = (item) => {

        return (
            <Grid container direction='column'>
                <Grid item onClick={btnClick.bind(this, item.type, item.category)} style={{cursor: 'pointer'}}>
                    <Grid container>
                        <Grid item>
                            <Grid container alignItems='center'>
                                <img style={{filter: "invert(100%)", opacity: "50%", marginRight: "10px"}} src={"data:image/png;base64," + item.image}/>

                                <Typography style={{
                                    fontSize: '14pt',
                                    overflow: 'hidden',
                                    height: '20pt',
                                    width: '200px',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                }}>{item.display_name}</Typography>
                            </Grid>
                        </Grid>
                        <Grid item style={{marginTop: '10px'}}>
                            {displayUnitValue(item)}
                        </Grid>
                    </Grid>
                </Grid>
            </Grid>
        )
    }

    return (
        <Grid container direction="column">
            <Grid item>
                <Grid container alignItems="center" style={{backgroundColor: "#222", borderRadius: "10px", padding: "5px"}}>
                    <Tooltip title="Back" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                        <IconButton onClick={()=>{
                            window.history.back()
                        }} variant="outlined" color="primary">
                            <ArrowBackIosOutlinedIcon fontSize="large"  style={{borderBottom: "0px solid #00A0E3"}}/>
                        </IconButton>
                    </Tooltip>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container direction="column">
                    <Grid item>
                        <Grid container direction="row">
                            {unitTypes !== undefined && unitTypes.types !== undefined ? unitTypes.types.map((item) => (
                                <Grid item
                                      key={"gazer-page_units-unittype-" + item.type}
                                      button
                                      onMouseEnter={() => handleEnter(item.type)}
                                      onMouseLeave={() => handleLeave(item.type)}
                                      style={btnStyle(item.type)}
                                >
                                    {displayItem(item)}
                                </Grid>
                            )) : <div/>}
                        </Grid>
                    </Grid>
                </Grid>

            </Grid>
        </Grid>
    );
}

export default PageAddUnit;
