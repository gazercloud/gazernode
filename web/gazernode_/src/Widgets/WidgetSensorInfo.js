import React, {useState} from "react";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";

export default function WidgetSensorInfo(props) {
    const [sensorId, setSensorID] = useState("")
    const [sensorType, setSensorType] = useState("")
    const [sensorName, setSensorName] = useState("")
    const [requestSensorPropertiesProcessing, setRequestSensorPropertiesProcessing] = useState(false)
    const requestSensorProperties = (paths) => {
        if (requestSensorPropertiesProcessing) {
            return
        }
        setRequestSensorPropertiesProcessing(true)

        let req = {
            "func": "host_read_items",
            "paths": paths,
            "layer": ""
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            for (let i = 0; i < result.items.length; i++) {
                                if (result.items[i].path.includes("/type"))
                                    setSensorType(result.items[i].value)
                                if (result.items[i].path.includes("/name"))
                                    setSensorName(result.items[i].value)
                            }

                            setRequestSensorPropertiesProcessing(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                setRequestSensorPropertiesProcessing(false)
            })
            .catch((err) => {
                setRequestSensorPropertiesProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }

    const requestSensorStart = () => {
        let req = {
            "func": "host_sensor_start",
            "id": sensorId
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                setRequestSensorPropertiesProcessing(false)
            })
            .catch((err) => {
                setRequestSensorPropertiesProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }

    const requestSensorStop = () => {
        let req = {
            "func": "host_sensor_stop",
            "id": sensorId
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(

                    );
                } else {
                    res.json().then(
                        (result) => {
                        }
                    );
                }
                setRequestSensorPropertiesProcessing(false)
            })
            .catch((err) => {
                setRequestSensorPropertiesProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering || sensorId !== props.SensorID) {
        setSensorID(props.SensorID)

        let paths = []
        paths.push("/sensors/" + props.SensorID + "/type")
        paths.push("/sensors/" + props.SensorID + "/name")
        requestSensorProperties(paths)

        setFirstRendering(false)
    }


    return (
        <div>
            <Typography style={{fontSize: "14pt"}}>{sensorName}</Typography>
            <Typography style={{fontSize: "10pt"}}>{sensorType}</Typography>
            <Button onClick={() => {requestSensorStart()}}>Start</Button>
            <Button onClick={() => {requestSensorStop()}}>Stop</Button>
        </div>
    )

}
