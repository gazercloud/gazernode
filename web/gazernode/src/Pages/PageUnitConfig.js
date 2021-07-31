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
import {Tooltip} from "@material-ui/core";
import Zoom from "@material-ui/core/Zoom";
import IconButton from "@material-ui/core/IconButton";
import ArrowBackIosOutlinedIcon from "@material-ui/icons/ArrowBackIosOutlined";
import Divider from "@material-ui/core/Divider";
import EditOutlinedIcon from "@material-ui/icons/EditOutlined";
import SaveIcon from '@material-ui/icons/Save';
import WidgetError from "../Widgets/WidgetError";

function PageUnitConfig(props) {
    const [unitConfig, setUnitConfig] = React.useState({})
    const [unitConfigMeta, setUnitConfigMeta] = React.useState([])
    const [unitId, setUnitId] = React.useState("")
    const [unitName, setUnitName] = React.useState("")

    const [configLoading, setConfigLoading] = React.useState(false)
    const [configLoaded, setConfigLoaded] = React.useState(false)

    const [configMetaLoading, setConfigMetaLoading] = React.useState(false)
    const [configMetaLoaded, setConfigMetaLoaded] = React.useState(false)


    const [messageError, setMessageError] = React.useState("")

    const requestSensorConfig = (unitId) => {
        if (configLoading || configLoaded) {
            return
        }

        setConfigLoading(true)

        let req = {
            id: unitId
        }
        Request('unit_get_config', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitId(result.id)
                            setUnitName(result.name)

                            let meta = JSON.parse(result.config_meta)
                            setUnitConfigMeta(meta)

                            let conf = JSON.parse(result.config)
                            setUnitConfig(conf)

                            console.log("UnitConfig", result)
                            setConfigLoaded(true)
                            setConfigLoading(false)
                            setConfigMetaLoaded(true)
                            setMessageError("")
                        }
                    );
                } else {
                    if (res.status !== 500) {
                        throw "123"
                    }

                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                            setConfigLoading(false)
                            RequestFailed()
                        }
                    );
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
            setConfigLoading(false)
        });
    }

    const requestSensorConfigMeta = (unitType) => {
        if (configMetaLoading || configMetaLoaded) {
            return
        }
        setConfigMetaLoading(true)

        let req = {
            type: unitType
        }
        Request('unit_type_config_meta', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitId("")
                            setUnitName("")

                            let meta = JSON.parse(result.config_meta)

                            setUnitConfigMeta(meta)
                            setUnitConfig({})
                            setConfigLoaded(true)
                            setConfigMetaLoaded(true)
                            setConfigMetaLoading(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setMessageError(result.error)
                            setConfigMetaLoading(false)
                            RequestFailed()
                        }
                    );
                }
            }).catch(res => {
            setMessageError(res.message)
            RequestFailed()
            setConfigMetaLoading(false)
        });
    }

    const requestSaveConfig = () => {
        let configAsString = JSON.stringify(unitConfig)

        if (unitId !== "") {
            let req = {
                id: unitId,
                name: unitName,
                config: configAsString
            }
            Request('unit_set_config', req)
                .then((res) => {
                    if (res.status === 200) {
                        res.json().then(
                            (result) => {
                                window.history.back()
                            }
                        );
                    } else {
                        res.json().then(
                            (result) => {
                                setMessageError(result.error)
                            }
                        );
                    }
                });
        } else {
            const unitType = new Buffer(props.UnitType, 'hex').toString();
            let req = {
                type: unitType,
                name: unitName,
                config: configAsString
            }
            Request('unit_add', req)
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
                                setMessageError(result.error)
                            }
                        );
                    }
                });
        }
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        if (props.UnitId === undefined || props.UnitId === "") {
            const unitType = new Buffer(props.UnitType, 'hex').toString();
            requestSensorConfigMeta(unitType)
        } else {
            const unitId = new Buffer(props.UnitId, 'hex').toString();
            requestSensorConfig(unitId)
        }
        setFirstRendering(false)
    }

    useEffect(() => {
        //const unitId = new Buffer(props.UnitId, 'hex').toString();
        const timer = setInterval(() => {
            if (props.UnitId === undefined || props.UnitId === "") {
                const unitType = new Buffer(props.UnitType, 'hex').toString();
                requestSensorConfigMeta(unitType)
            } else {
                const unitId = new Buffer(props.UnitId, 'hex').toString();
                requestSensorConfig(unitId)
            }
        }, 1000);
        return () => clearInterval(timer);
    });

    const displayErrorMessage = () => {
        return (
            <div style={{color: "#F00", fontSize: "14pt"}}>{messageError}</div>
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

                    <Divider orientation="vertical" flexItem style={{marginLeft: "20px", marginRight: "20px"}} />

                    <Tooltip title="Edit" TransitionComponent={Zoom} style={{border: "1px solid #333", marginRight: "5px"}}>
                        <IconButton onClick={requestSaveConfig.bind(this)} variant="outlined" color="primary" disabled={!configLoaded || !configMetaLoaded}>
                            <SaveIcon fontSize="large"  style={{borderBottom: "0px solid #00A0E3"}}/>
                        </IconButton>
                    </Tooltip>
                </Grid>

            </Grid>
            <Grid item>
                <Grid container alignItems="center" spacing={1} style={{marginTop: "20px"}}>
                    <Grid item style={{minWidth: "100px", color: "#555"}}>Id:</Grid>
                    <Grid item style={{color: "#555"}}>
                        {unitId}
                    </Grid>
                </Grid>
            </Grid>
            <Grid item>
                <Grid container alignItems="center" spacing={1} style={{marginBottom: "20px"}}>
                    <Grid item style={{minWidth: "100px"}}>Name:</Grid>
                    <Grid item>
                        <TextField
                            disabled={!configLoaded || !configMetaLoaded}
                            value={unitName}
                            placeholder={"Name of the unit"}
                            onChange={(ev)=>{
                                setUnitName(ev.target.value)
                            }}/>
                    </Grid>
                </Grid>
            </Grid>

            <Grid item>
                <ConfigSensor ConfigMeta={unitConfigMeta} Config={unitConfig} OnConfigChanged={(v) => setUnitConfig(v)}/>
            </Grid>

            <Grid>
                {displayErrorMessage()}
            </Grid>

            <div>
                <WidgetError Message={messageError} />
            </div>

        </Grid>
    );
}

export default PageUnitConfig;
