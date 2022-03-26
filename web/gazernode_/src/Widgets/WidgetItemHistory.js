import React, {useEffect, useRef, useState} from 'react';
import Request, {RequestFailed} from "../request";
import ButtonGroup from "@material-ui/core/ButtonGroup";
import Button from "@material-ui/core/Button";
import WidgetTimeChart from "./WidgetTimeChart/WidgetTimeChart";
import JSZip from "jszip";
import {LinearProgress} from "@material-ui/core";
import NewTimeChart from "./WidgetTimeChart/TimeChart";

export default function WidgetItemHistory(props) {
    const [dataItemHistory, setDataItemHistory] = React.useState([])
    const [timeRangeLast, setTimeRangeLast] = React.useState(5)
    const [lastUpdateTime, setLastUpdateTime] = React.useState(0)
    const [stateLoading, setStateLoading] = React.useState(false)
    const [renderCounter, setRenderCounter] = React.useState(0)

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestDataItemHistory(props.DataItemName, minMaxTime(timeRangeLast), false)
            setRenderCounter(renderCounter + 1)
        }, 1000);
        return () => clearInterval(timer);
    });

    const alignGroupTimeRange = (groupTimeRange) => {
        if (groupTimeRange < 100000) {
            groupTimeRange = 1
        }

        //console.log("groupTimeRange", groupTimeRange)

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

        //console.log("groupTimeRange FINAL:", groupTimeRange)

        return groupTimeRange
    }

    const requestDataItemHistory = (dataItemName, minmax, force) => {
        if (stateLoading) {
            return
        }

        let updateTimeout = 10000
        if (timeRangeLast <= 5) {
            updateTimeout = 1000
        }

        if (!force) {
            if (((new Date().getTime()) - lastUpdateTime) <= updateTimeout) {
                return
            }
        }

        //let minmax =

        let req = {
            name: dataItemName,
            dt_begin: minmax.min,
            dt_end: minmax.max,
            group_time_range: minmax.groupTimeRange,
            out_format: "zip"
        }

        setStateLoading(true)

        Request('data_item_history_chart', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setLastUpdateTime(new Date().getTime())
                            const byteCharacters = atob(result.data);
                            let new_zip = new JSZip();
                            new_zip.loadAsync(byteCharacters)
                                .then(function (zip) {
                                    zip.file("data").async("string").then((v) => {
                                        let obj = JSON.parse(v)

                                        let items = [];
                                        items = obj.items
                                        let newData = items
                                        setDataItemHistory(newData)

                                    });
                                });
                        }
                    );

                    setStateLoading(false)
                    return
                }

                setStateLoading(false)
                RequestFailed()

            }).catch(res => {
            setStateLoading(false)
            RequestFailed()
        });
    }

    const width = () => {
        let w = window.innerWidth
        if (window.innerWidth > 600) {
            w -= 310
        } else {
            w -= 65
        }
        return w
    }

    const minMaxTime = (timeRangeLastValue) => {
        let w = width() // width
        let r = timeRangeLastValue * 60.0 * 1000000.0
        let timePerPixel = r / w
        //console.log("timePerPixel", timePerPixel)

        let groupTimeRange = alignGroupTimeRange(timePerPixel)

        let dtNow = new Date().getTime() * 1000
        dtNow = dtNow - (dtNow % groupTimeRange)

        let res = {max: dtNow, min: dtNow - timeRangeLastValue * 60000000, groupTimeRange: groupTimeRange}
        return res
    }

    const renderChart = () => {
        return (
            <WidgetTimeChart Data={dataItemHistory} ChartWidth={width()} MinTime={minMaxTime(timeRangeLast).min} MaxTime={minMaxTime(timeRangeLast).max}>
            </WidgetTimeChart>
        );
    }

    const setLastRange = (g, v) => {
        setStateLoading(false)
        setTimeRangeLast(g)
        setLastUpdateTime(0)
        setDataItemHistory([])
        requestDataItemHistory(props.DataItemName, minMaxTime(g), true)
    }

    const btnStyle = (v) => {
        if (v === timeRangeLast) {
            return {
                color: "#0f0",
                minWidth: "100px"
            }
        }
        return {
            color: "#00A0E3",
            minWidth: "100px"
        }
    }

    return (
        <div>
            <div>
                <ButtonGroup color="primary" aria-label="outlined primary button group">
                    <Button onClick={setLastRange.bind(this, 5)} style={btnStyle(5)}>5m</Button>
                    <Button onClick={setLastRange.bind(this, 60)} style={btnStyle(60)}>60m</Button>
                    <Button onClick={setLastRange.bind(this, 60 * 24)} style={btnStyle(60 * 24)}>24H</Button>
                </ButtonGroup>
            </div>
            <div style={{marginTop: "10px"}}>{renderChart()}</div>
            <div>{stateLoading ? <div><LinearProgress /></div> : <div><LinearProgress value={0} variant="determinate" /></div>}</div>
        </div>
    );

}

