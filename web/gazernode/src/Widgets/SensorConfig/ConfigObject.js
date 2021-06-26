import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import ConfigText from "./ConfigText";
import ConfigNumber from "./ConfigNumber";
import ConfigTable from "./ConfigTable";
import Grid from "@material-ui/core/Grid";
import ConfigString from "./ConfigString";

export default function ConfigObject(props) {

    const drawItem = (meta, data) => {
        console.log("ConfigObject drawItem", meta)
        if (meta.type === "string") {
            console.log("ConfigObject drawItem text", meta)
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

    console.log("ConfigObject render", props.Meta)
    return (

        <Paper variant="outlined" style={{padding: "5px"}}>
            {props.Meta.map((item) => (
                <div style={{margin: "5px"}}>
                    {drawItem(item, props.Data)}
                </div>
            ))
            }
        </Paper>
    )
}
