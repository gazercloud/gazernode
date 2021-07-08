import React, {useEffect, useState} from 'react';
import Paper from "@material-ui/core/Paper";
import WidgetDataItems from "./WidgetDataItems";
import WidgetDataItemDetails from "./WidgetDataItemDetails";
import Grid from "@material-ui/core/Grid";
import WidgetSensorInfo from "./WidgetSensorInfo";
import Request from "../request";
import {CartesianGrid, Legend, Line, LineChart, XAxis, YAxis} from "recharts";
import ButtonGroup from "@material-ui/core/ButtonGroup";
import Button from "@material-ui/core/Button";
import WidgetTimeChart from "./WidgetTimeChart/WidgetTimeChart";

export default function WidgetItemHistory(props) {
    const [dataItemHistory, setDataItemHistory] = React.useState([])
    const [timeRangeLast, setTimeRangeLast] = React.useState(5)
    const [lastUpdateTime, setLastUpdateTime] = React.useState(0)
    const [requesting, setRequesting] = React.useState(false)
    const [renderCounter, setRenderCounter] = React.useState(0)

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestDataItemHistory(props.DataItemName)
            setRenderCounter(renderCounter + 1)
        }, 1000);
        return () => clearInterval(timer);
    });

    const alignGroupTimeRange = (groupTimeRange) => {
        if (groupTimeRange < 100000) {
            groupTimeRange = 1
        }

        console.log("groupTimeRange", groupTimeRange)

        if (groupTimeRange >= 100000 && groupTimeRange < 200000) {
            groupTimeRange = 100000 // By 0.2 sec
        }

        if (groupTimeRange >= 200000 && groupTimeRange < 500000) {
            groupTimeRange = 200000 // By 0.2 sec
        }

        if (groupTimeRange >= 500000 && groupTimeRange < 1000000) {
            groupTimeRange = 500000 // By 0.5 sec
        }

        if (groupTimeRange >= 1000000 && groupTimeRange < 5 * 1000000) {
            groupTimeRange = 1000000 // By 1 sec
        }

        if (groupTimeRange >= 5 * 1000000 && groupTimeRange < 15 * 1000000) {
            groupTimeRange = 5 * 1000000 // By 5 sec
        }

        if (groupTimeRange >= 15 * 1000000 && groupTimeRange < 30 * 1000000) {
            groupTimeRange = 15 * 1000000 // By 15 sec
        }

        if (groupTimeRange >= 30 * 1000000 && groupTimeRange < 60 * 1000000) {
            groupTimeRange = 30 * 1000000 // By 30 sec
        }

        if (groupTimeRange >= 60 * 1000000 && groupTimeRange < 2 * 60 * 1000000) {
            groupTimeRange = 60 * 1000000 // By minute
        }

        if (groupTimeRange >= 2 * 60 * 1000000 && groupTimeRange < 3 * 60 * 1000000) {
            groupTimeRange = 2 * 60 * 1000000 // By 2 minute
        }

        if (groupTimeRange >= 3 * 60 * 1000000 && groupTimeRange < 4 * 60 * 1000000) {
            groupTimeRange = 3 * 60 * 1000000 // By 3 minute
        }

        if (groupTimeRange >= 4 * 60 * 1000000 && groupTimeRange < 5 * 60 * 1000000) {
            groupTimeRange = 4 * 60 * 1000000 // By 4 minute
        }

        if (groupTimeRange >= 5 * 60 * 1000000 && groupTimeRange < 10 * 60 * 1000000) {
            groupTimeRange = 5 * 60 * 1000000 // By 5 minute
        }

        if (groupTimeRange >= 10 * 60 * 1000000 && groupTimeRange < 20 * 60 * 1000000) {
            groupTimeRange = 10 * 60 * 1000000 // By 10 minute
        }

        if (groupTimeRange >= 20 * 60 * 1000000 && groupTimeRange < 30 * 60 * 1000000) {
            groupTimeRange = 20 * 60 * 1000000 // By 20 minute
        }

        if (groupTimeRange >= 30 * 60 * 1000000 && groupTimeRange < 60 * 60 * 1000000) {
            groupTimeRange = 30 * 60 * 1000000 // By 30 minute
        }

        if (groupTimeRange >= 60 * 60 * 1000000 && groupTimeRange < 3 * 60 * 60 * 1000000) {
            groupTimeRange = 60 * 60 * 1000000 // By 60 minutes
        }

        if (groupTimeRange >= 3 * 60 * 60 * 1000000 && groupTimeRange < 6 * 60 * 60 * 1000000) {
            groupTimeRange = 3 * 60 * 60 * 1000000 // By 3 Hours
        }

        if (groupTimeRange >= 6 * 60 * 60 * 1000000 && groupTimeRange < 12 * 60 * 60 * 1000000) {
            groupTimeRange = 6 * 60 * 60 * 1000000 // By 6 Hours
        }

        if (groupTimeRange >= 12 * 60 * 60 * 1000000 && groupTimeRange < 24 * 60 * 60 * 1000000) {
            groupTimeRange = 12 * 60 * 60 * 1000000 // By 6 Hours
        }

        if (groupTimeRange >= 24 * 60 * 60 * 1000000) {
            groupTimeRange = 24 * 60 * 60 * 1000000 // By day
        }

        console.log("groupTimeRange FINAL:", groupTimeRange)

        return groupTimeRange
    }

    const requestDataItemHistory = (dataItemName) => {
        if (requesting) {
            return
        }
        setRequesting(true)

        /*let w = width() // width
        let r = timeRangeLast * 60.0 * 1000000.0
        let timePerPixel = r / w*/

        let updateTimeout = 10000
        if (timeRangeLast <= 5) {
            updateTimeout = 1000
        }

        if (((new Date().getTime()) - lastUpdateTime) <= updateTimeout) {
            setRequesting(false)
            return
        }

        //console.log("timePerPixel", timePerPixel)

        /*let groupTimeRange = alignGroupTimeRange(timePerPixel)

        let dtNow = new Date().getTime() * 1000
        dtNow = dtNow - (dtNow % groupTimeRange)*/

        let minmax = minMaxTime()

        let req = {
            name: dataItemName,
            dt_begin: minmax.min,
            dt_end: minmax.max,
            group_time_range: minmax.groupTimeRange
        }
        Request('data_item_history_chart', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setLastUpdateTime(new Date().getTime())
                            let items = [];
                            items = result.items
                            let newData = items
                            setDataItemHistory(newData)
                            setRequesting(false)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            setRequesting(false)
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const width = () => {
        let w = window.innerWidth - 300
        if (w < 300)
            return 300
        return w
    }

    const minMaxTime = () => {
        let w = width() // width
        let r = timeRangeLast * 60.0 * 1000000.0
        let timePerPixel = r / w
        console.log("timePerPixel", timePerPixel)

        let groupTimeRange = alignGroupTimeRange(timePerPixel)

        let dtNow = new Date().getTime() * 1000
        dtNow = dtNow - (dtNow % groupTimeRange)

        let res = {max: dtNow, min: dtNow - timeRangeLast * 60000000, groupTimeRange: groupTimeRange}
        return res
    }

    const renderChart = () => {
        return (
            <WidgetTimeChart Data={dataItemHistory} ChartWidth={width()} MinTime={minMaxTime().min} MaxTime={minMaxTime().max}>
            </WidgetTimeChart>
        );
    }

    const setLastRange = (g, v) => {
        setRequesting(false)
        setTimeRangeLast(g)
        setLastUpdateTime(0)
        setDataItemHistory([])
    }

    const btnStyle = (v) => {
        if (v === timeRangeLast) {
            return {
                color: "#0f0",
                minWidth: "100px"
            }
        }
        return {
            color: "#f00",
            minWidth: "100px"
        }
    }

    return (
        <div>
            <div>
                <div>
                    <ButtonGroup color="primary" aria-label="outlined primary button group">
                        <Button onClick={setLastRange.bind(this, 1)} style={btnStyle(1)}>1m</Button>
                        <Button onClick={setLastRange.bind(this, 5)} style={btnStyle(5)}>5m</Button>
                        <Button onClick={setLastRange.bind(this, 30)} style={btnStyle(30)}>30m</Button>
                        <Button onClick={setLastRange.bind(this, 60)} style={btnStyle(60)}>60m</Button>
                    </ButtonGroup>
                </div>
                <div>
                    <ButtonGroup color="primary" aria-label="outlined primary button group">
                        <Button onClick={setLastRange.bind(this, 60 * 3)} style={btnStyle(60 * 3)}>3H</Button>
                        <Button onClick={setLastRange.bind(this, 60 * 6)} style={btnStyle(60 * 6)}>6H</Button>
                        <Button onClick={setLastRange.bind(this, 60 * 12)} style={btnStyle(60 * 12)}>12H</Button>
                        <Button onClick={setLastRange.bind(this, 60 * 24)} style={btnStyle(60 * 24)}>24H</Button>
                    </ButtonGroup>
                </div>
                <div>
                    <ButtonGroup color="primary" aria-label="outlined primary button group">
                        <Button onClick={setLastRange.bind(this, 7 * 60 * 24)} style={btnStyle(7 * 60 * 24)}>7
                            days</Button>
                        <Button onClick={setLastRange.bind(this, 30 * 60 * 24)} style={btnStyle(30 * 60 * 24)}>30
                            days</Button>
                        <Button onClick={setLastRange.bind(this, 90 * 60 * 24)} style={btnStyle(90 * 60 * 24)}>90
                            days</Button>
                        <Button onClick={setLastRange.bind(this, 180 * 60 * 24)} style={btnStyle(180 * 60 * 24)}>180
                            days</Button>
                    </ButtonGroup>
                </div>
            </div>
            <div style={{marginTop: "10px"}}>{renderChart()}</div>
            <div>{requesting ? <div>Requesting</div> : <div>...</div>}</div>
        </div>
    );

}

