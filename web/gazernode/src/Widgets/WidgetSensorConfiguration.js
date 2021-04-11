import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Paper from "@material-ui/core/Paper";
import Request from "../request";
import ConfigSensor from "./SensorConfig/ConfigSensor";
import Button from "@material-ui/core/Button";

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    values: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    expand: {
        transform: 'rotate(0deg)',
        marginLeft: 'auto',
        transition: theme.transitions.create('transform', {
            duration: theme.transitions.duration.shortest,
        }),
    },
    expandOpen: {
        transform: 'rotate(180deg)',
    },
}));

function WidgetSensorsConfiguration(props) {
    const classes = useStyles();

    const [configLoaded, setConfigLoaded] = React.useState(false)
    const [configMetaLoaded, setConfigMetaLoaded] = React.useState(false)

    const [sensorId, setSensorId] = React.useState("")
    const [sensorType, setSensorType] = React.useState("")
    const [configMeta, setConfigMeta] = React.useState([])
    const [config, setConfig] = React.useState([])

    const [configOriginal, setConfigOriginal] = React.useState("")

    const [errorMessage, setErrorMessage] = React.useState("")

    const loadSensorConfig = (sensorId) => {
        console.log("loadSensorConfig", sensorId)
        let req = {
            "func": "host_get_sensor_config",
            "id": sensorId
        }

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    res.text().then(
                        (result) => {
                            try {
                                let obj = JSON.parse(result);
                                console.log("loadSensorConfig ok")
                                console.log(obj)
                                setConfig(obj)
                            } catch (e) {
                                setErrorMessage("loadSensorConfig Wrong json")
                                setConfig({})
                                console.log("loadSensorConfig Wrong json", e)
                                console.log(result)
                            }

                        }
                    )
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
    }


    const loadSensorConfigMeta = (sType) => {
        console.log("loadSensorConfigMeta", sType)
        let req = {
            "func": "host_get_sensor_config_meta",
            "sensor_type": sType
        }

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    res.text().then(
                        (result) => {
                            try {
                                let obj = JSON.parse(result);
                                console.log("loadSensorConfigMeta ok")
                                console.log(obj)
                                setConfigMeta(obj)
                                loadSensorConfig(props.SensorID)
                            } catch (e) {
                                setErrorMessage("Wrong json")
                                console.log("Wrong json", e)
                                console.log(result)
                            }

                        }
                    )
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
    }


    const loadSensorType = (sensorId) => {
        console.log("loadSensorType", sensorId)
        let req = {
            "func": "host_get_sensor_type",
            "id": sensorId
        }

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            console.log("loadSensorType ok")
                            console.log(result)
                            loadSensorConfigMeta(result.sensor_type)
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

    const saveConfig = () => {
        console.log("saveConfig")
        let req = {
            "func": "host_set_sensor_config",
            "id": sensorId,
            "config": JSON.stringify(config)
        }

        Request('/api/request', req)
            .then((res) => {
                if (res.status === 200) {
                } else {
                    res.json().then(
                        (result) => {
                            setErrorMessage(result.error)
                        }
                    );
                }
            });
    }


    if (props.SensorID !== sensorId) {
        setSensorId(props.SensorID)

        setSensorType("")
        setConfig(undefined)
        setConfigMeta(undefined)
        setConfigLoaded(false)
        setConfigMetaLoaded(false)
        loadSensorType(props.SensorID)
    }

    const hasChanges = () => {
        if (config !== configOriginal)
            return true
        return false
    }

    return (
        <Paper>
            <ConfigSensor Config={config} ConfigMeta={configMeta} OnConfigChanged={(c)=>{
                setConfig(c)
            }}/>
            <Button onClick={()=>{
                saveConfig()
            }}>
                Save
            </Button>
        </Paper>
    );
}

export default WidgetSensorsConfiguration;
