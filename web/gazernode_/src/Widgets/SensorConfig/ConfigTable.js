import * as React from "react";
import {useState} from "react";
import Paper from "@material-ui/core/Paper";
import TextField from "@material-ui/core/TextField";
import Table from "@material-ui/core/Table";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import Grid from "@material-ui/core/Grid";
import Button from "@material-ui/core/Button";
import TableBody from "@material-ui/core/TableBody";
import ConfigObject from "./ConfigObject";
import ListItem from "@material-ui/core/ListItem";
import Typography from "@material-ui/core/Typography";
import {TableContainer} from "@material-ui/core";

export default function ConfigTable(props) {

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (ev, key) => {
        setCurrentItem(ev)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }
    const btnStyle = (key) => {
        if (currentItem === key) {
            return {
                cursor: "pointer",
                backgroundColor: "#52bdff",
            }
        } else {
            if (hoverItem === key) {
                return {
                    cursor: "pointer",
                    backgroundColor: "#b8e4ff"
                }
            } else {
                return {
                    cursor: "pointer",
                    backgroundColor: "#222"
                }
            }
        }
    }

    const getDataOfItem = (index) => {
        if (props.Data === undefined)
            return undefined
        return props.Data[index]
    }

    const filterInt = (value) => {
        if (/^(\-|\+)?([0-9]+|Infinity)$/.test(value))
            return Number(value);
        return NaN;
    }

    const limitString = (str, maxSize) => {
        if (str.length > maxSize) {
            return str.substr(0, maxSize) + "..."
        }
        return str
    }

    return (
        <Grid container={true} direction="column">
            <Grid item>
                <Typography>{props.Meta.display_name}:</Typography>
            </Grid>
            <Grid item>
                <Grid container={true} direction="row" justify="flex-start" alignItems="center" style={{marginTop: "10px", marginBottom: "10px"}}>
                <Button
                    color="primary"
                    variant="outlined"
                    style={{minWidth: "120px", marginRight: "10px"}}
                    onClick={() => {
                        let workCopy = []
                        if (props.Data !== undefined)
                            workCopy = JSON.parse(JSON.stringify(props.Data));
                        let defObject = {}
                        for (let i = 0; i < props.Meta.children.length; i++) {
                            let defValue = undefined
                            if (props.Meta.children[i].default_value !== undefined) {
                                defValue = props.Meta.children[i].default_value

                                if (props.Meta.children[i].type === "bool") {
                                    if (defValue === "true" || defValue === "1") {
                                        defValue = true
                                    } else {
                                        defValue = false
                                    }
                                }
                                if (props.Meta.children[i].type === "num")
                                    defValue = parseFloat(defValue)
                                if (props.Meta.children[i].type === "table")
                                    defValue = []
                                if (props.Meta.children[i].type === "object")
                                    defValue = {}

                            } else {
                                if (props.Meta.children[i].type === "bool")
                                    defValue = false
                                if (props.Meta.children[i].type === "num")
                                    defValue = 0
                                if (props.Meta.children[i].type === "string")
                                    defValue = ""
                                if (props.Meta.children[i].type === "text")
                                    defValue = ""
                                if (props.Meta.children[i].type === "table")
                                    defValue = []
                                if (props.Meta.children[i].type === "object")
                                    defValue = {}
                            }
                            defObject[props.Meta.children[i].name] = defValue
                        }
                        workCopy.push(defObject)
                        props.OnChangedValue(props.Meta.name, workCopy)
                    }}>Add row</Button>

                <Button
                    color="primary"
                    variant="outlined"
                    onClick={() => {
                        let currentItemAsInt = filterInt(currentItem);
                        if (isNaN(currentItemAsInt))
                            return;
                        let workCopy = []
                        if (props.Data !== undefined)
                            workCopy = JSON.parse(JSON.stringify(props.Data));
                        workCopy.splice(currentItemAsInt, 1)
                        props.OnChangedValue(props.Meta.name, workCopy)
                    }}>Remove row</Button>
                </Grid>
            </Grid>

            <Grid item>
                <Grid container>
                    <Grid item>
                        <TableContainer style={{maxWidth: "600px", maxHeight: "400px"}}>
                        <Table size="small" style={{minWidth: "400px"}}>

                            <TableBody>
                                <TableRow style={{backgroundColor: "#222233"}}>
                                    {props.Meta.children.map((item) => (
                                        <TableCell style={{maxWidth: "100px"}}>
                                            <div>{item.display_name}</div>
                                        </TableCell>
                                    ))
                                    }
                                </TableRow>
                                {props.Data !== undefined ? props.Data.map((dataItem, dataItemIndex) =>
                                    (
                                        <TableRow
                                            onMouseEnter={() => handleEnter(dataItemIndex)}
                                            onMouseLeave={() => handleLeave(dataItemIndex)}
                                            style={btnStyle(dataItemIndex)}
                                            onClick={btnClick.bind(this, dataItemIndex)}
                                        >
                                            {props.Meta.children.map((item) => (
                                                <TableCell  style={{maxWidth: "100px"}}>
                                                    {(item.type !== "table" && item.type !== "object") ? limitString(dataItem[item.name], 10) : "obj"}
                                                </TableCell>
                                            ))}
                                        </TableRow>
                                    )) : <div/>
                                }
                            </TableBody>
                        </Table>
                        </TableContainer>
                    </Grid>
                    <Grid item>
                        <div style={{minWidth: "100px", minHeight: "100px", margin: "5px"}}>
                            {props.Data === undefined || props.Data.constructor !== Array || currentItem === "" || currentItem >= props.Data.length || currentItem < 0 ?
                                <div>No row selected</div>
                                :
                                <div>
                                    <ConfigObject Meta={props.Meta.children} Data={getDataOfItem(currentItem)}
                                                  OnChangedValue={
                                                      (n, v) => {
                                                          let workCopy = JSON.parse(JSON.stringify(props.Data));
                                                          workCopy[currentItem] = v
                                                          props.OnChangedValue(props.Meta.name, workCopy)
                                                      }
                                                  }/>
                                </div>

                            }
                        </div>
                    </Grid>
                </Grid>
            </Grid>
        </Grid>

    )
}
