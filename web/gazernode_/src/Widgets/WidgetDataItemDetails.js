import React, {useEffect, useState} from 'react';
import Typography from "@material-ui/core/Typography";
import WidgetDataItemHistory from "./WidgetDataItemHistory";
import WidgetSensorInfo from "./WidgetSensorInfo";

export default function WidgetDataItemDetails(props) {
    const [currentPath, setCurrentPath] = React.useState(props.Path)
    const [value, setValue] = React.useState({})
    const [lastPath, setLastPath] = useState("")
    const formatDateTime = (date) => {
        let yyyy = date.getFullYear();
        let mm = date.getMonth() + 1;
        if (mm < 10) mm = '0' + mm;
        let dd = date.getDate();
        if (dd < 10) dd = '0' + dd;
        let hh = date.getHours();
        if (hh < 10) hh = '0' + hh;
        let min = date.getMinutes();
        if (min < 10) min = '0' + min;
        let ss = date.getSeconds();
        if (ss < 10) ss = '0' + ss;
        let fff = date.getMilliseconds();
        if (fff < 10) {
            fff = '00' + fff;
        } else {
            if (fff < 100)
                fff = '0' + fff
        }

        return yyyy + '-' + mm + '-' + dd + ' ' + hh + ':' + min + ':' + ss + '.' + fff;
    }


    const requestItem = (path) => {
        console.log("requestItem")
        let req = {
            "func": "host_children",
            "path": path
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {

                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                //setProcessingLoadParameter(false)
            })
            .catch((err) => {
                //setProcessingLoadParameter(false)
                //setErrorMessage("Unknown error")
            })
    }

    const [requestValuesProcessing, setRequestValuesProcessing] = React.useState(false)
    const requestValues = (path) => {
        return
        if (requestValuesProcessing) {
            return
        }
        setRequestValuesProcessing(true)

        let req = {
            "func": "host_read_items",
            "paths": [path]
        }
        fetch("/api/request?request=" + JSON.stringify(req))
            .then(res => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            for (let i = 0; i < result.items.length; i++) {
                                let dt1 = new Date(result.items[i].dt)
                                let dt2 = ""
                                if (dt1 !== undefined) {
                                    dt2 = formatDateTime(dt1)
                                }

                                let val = {
                                    "value":result.items[i].value,
                                    "dt":dt2,
                                    "q":result.items[i].q,
                                    "uom":result.items[i].uom
                                }
                                setValue(val)
                            }
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
                setRequestValuesProcessing(false)
            })
            .catch((err) => {
                setRequestValuesProcessing(false)
                //setErrorMessage("Unknown error")
            })
    }

    const getValue = () => {
        return value.value
    }
    const getUOM = () => {
        return value.uom
    }

    const getDateTime = () => {
        return value.dt
    }

    const [firstRendering, setFirstRendering] = React.useState(true)
    if (firstRendering || lastPath !== props.Path ) {
        console.log("firstRendering DataItemDetails")
        console.log(props.Path)
        setLastPath(props.Path)
        setCurrentPath(props.Path)
        requestItem(props.Path)

        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestValues(currentPath)
        }, 500);
        return () => clearInterval(timer);
    });



    const drawDataItemName = () => {

        if (currentPath.indexOf("/sensors/") === 0)
        {
            let indexOfSensorIdEnd = currentPath.indexOf("/", 9)
            if (indexOfSensorIdEnd > 10) {
                let indexOfDataWordEnd = currentPath.indexOf("/", indexOfSensorIdEnd + 1)
                let path = currentPath.substring(indexOfSensorIdEnd + 1)
                if (indexOfDataWordEnd > 0) {
                    path = currentPath.substring(indexOfDataWordEnd + 1)
                }

                let sId = currentPath.substring(9, indexOfSensorIdEnd)
                return (
                    <div>
                        {props.ShowSensorInfo ? <WidgetSensorInfo SensorID={sId}/> : <div/>}
                        <Typography style={{fontSize: "16pt"}}>{path}</Typography>
                    </div>
                )
            }
        }

        return (
            <div>
                {currentPath}
            </div>
        )
    }

    if (currentPath === "" || currentPath === undefined || currentPath === null) {
        return (
            <div>
                No item to show
            </div>
        )
    }

    return (
        <div>
            <div>{drawDataItemName()}</div>
            <div>
                <span style={{fontSize: "24pt", marginTop: "20px", marginBottom: "20px"}}>{getValue()}</span>
                <span style={{fontSize: "16pt", marginTop: "20px", marginBottom: "20px", color: "#AAAAAA", marginLeft: "10px"}}>{getUOM()}</span>
            </div>
            <div>
                <span style={{fontSize: "16pt", marginTop: "20px", marginBottom: "20px", color: "#AAAAAA"}}>{getDateTime()}</span>
            </div>
            <div style={{maxWidth: "400px", minHeight: "600px"}}>
                <WidgetDataItemHistory Path={currentPath} />
            </div>
        </div>
    );
}


