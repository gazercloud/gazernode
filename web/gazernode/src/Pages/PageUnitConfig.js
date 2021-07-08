import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from '@material-ui/icons/Memory';
import Grid from "@material-ui/core/Grid";
import Request from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";
import ConfigSensor from "../Widgets/SensorConfig/ConfigSensor";
import TextField from "@material-ui/core/TextField";

function PageUnitConfig(props) {
    const [unitValues, setUnitValues] = React.useState([])
    const [unitConfig, setUnitConfig] = React.useState({})
    const [unitConfigMeta, setUnitConfigMeta] = React.useState([])
    const [unitId, setUnitId] = React.useState("")
    const [unitName, setUnitName] = React.useState("")
    const [errorMessage, setErrorMessage] = React.useState("")

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

    const requestSensorConfig = (unitId) => {
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
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const requestSensorConfigMeta = (unitType) => {
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
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
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
                                setErrorMessage(result.error)
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
                                setErrorMessage(result.error)
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
        }, 500);
        return () => clearInterval(timer);
    });

    const displayErrorMessage = () => {
        return (
            <div style={{color: "#F00", fontSize: "14pt"}}>{errorMessage}</div>
        )
    }

    return (
        <Grid container direction="column">
            <Grid item>
                <Button variant='outlined' color='primary' onClick={()=>{props.OnNavigate('#form=units')}}>
                    Back to the Units
                </Button>
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

            <Grid item>
                <Button variant='outlined' color='primary' style={{minWidth: '70px', margin: '5px'}}
                        onClick={requestSaveConfig.bind(this)}
                >SAVE</Button>
            </Grid>
            <Grid>
                {displayErrorMessage()}
            </Grid>

        </Grid>
    );
}

export default PageUnitConfig;
