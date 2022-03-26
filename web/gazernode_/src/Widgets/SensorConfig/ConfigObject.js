import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import ConfigText from "./ConfigText";
import ConfigNumber from "./ConfigNumber";
import ConfigTable from "./ConfigTable";
import Grid from "@material-ui/core/Grid";
import ConfigString from "./ConfigString";
import ConfigBool from "./ConfigBool";

export default function ConfigObject(props) {

    const drawItem = (meta, data) => {
        console.log("ConfigObject drawItem", meta)
        if (meta.type === "bool") {
            if (data[meta.name] === undefined) {
                if (meta.default_value !== undefined && meta.default_value !== "") {
                    if (meta.default_value === "true")
                        data[meta.name] = true
                    else
                        data[meta.name] = false
                } else {
                    data[meta.name] = ""
                }
            }
            return (
                <ConfigBool Meta={meta} Data={data !== undefined ? data[meta.name] : undefined}
                              OnChangedValue={
                                  (n, v) => {
                                      let workCopy = JSON.parse(JSON.stringify(props.Data));
                                      workCopy[n] = v
                                      props.OnChangedValue(meta.name, workCopy)
                                  }
                              }/>
            )
        }
        if (meta.type === "string") {
            console.log("ConfigObject drawItem text", meta)
            if (data[meta.name] === undefined) {
                if (meta.default_value !== undefined && meta.default_value !== "") {
                    data[meta.name] = meta.default_value
                } else {
                    data[meta.name] = ""
                }
            }
            return (
                <ConfigString Meta={meta} Data={data !== undefined ? data[meta.name] : undefined}
                            OnChangedValue={
                                (n, v) => {
                                    let workCopy = JSON.parse(JSON.stringify(props.Data));
                                    workCopy[n] = v
                                    props.OnChangedValue(meta.name, workCopy)
                                }
                            }/>
            )
        }
        if (meta.type === "text") {
            if (data[meta.name] === undefined) {
                if (meta.default_value !== undefined && meta.default_value !== "") {
                    data[meta.name] = meta.default_value
                } else {
                    data[meta.name] = ""
                }
            }

            console.log("ConfigObject drawItem text", meta)
            return (
                <ConfigText Meta={meta} Data={data !== undefined ? data[meta.name] : undefined}
                            OnChangedValue={
                                (n, v) => {
                                    let workCopy = JSON.parse(JSON.stringify(props.Data));
                                    workCopy[n] = v
                                    props.OnChangedValue(meta.name, workCopy)
                                }
                            }/>
            )
        }
        if (meta.type === "num") {
            if (data[meta.name] === undefined) {
                if (meta.default_value !== undefined && meta.default_value !== "") {
                    data[meta.name] = parseFloat(meta.default_value)
                } else {
                    data[meta.name] = 0
                }
            }

            return (
                <ConfigNumber Meta={meta} Data={data !== undefined ? data[meta.name] : undefined}
                              OnChangedValue={
                                  (n, v) => {
                                      let workCopy = JSON.parse(JSON.stringify(props.Data));
                                      workCopy[n] = v
                                      props.OnChangedValue(meta.name, workCopy)
                                  }
                              }/>
            )
        }

        if (meta.type === "table") {
            if (data[meta.name] === undefined) {
                data[meta.name] = []
            }
            return (
                <ConfigTable Meta={meta} Data={data !== undefined ? data[meta.name] : undefined}
                             OnChangedValue={
                                 (n, v) => {
                                     let workCopy = JSON.parse(JSON.stringify(props.Data));
                                     workCopy[n] = v
                                     props.OnChangedValue(meta.name, workCopy)
                                 }
                             }/>
            )
        }
    }

    if (props.Meta === undefined || props.Meta === null) {
        return (
            <Paper variant="outlined" style={{padding: "5px"}}>
            </Paper>
        )
    }



    return (

        <div>
            {props.Meta.map((item) => (
                <div>
                    {drawItem(item, props.Data)}
                </div>
            ))
            }
        </div>
    )
}
