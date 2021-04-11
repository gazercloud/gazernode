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
                    backgroundColor: "#FFFFFF"
                }
            }
        }
    }

    const getDataOfItem = (index) => {
        if (props.Data === undefined)
            return undefined
        return props.Data[index]
    }

    return (
        <Paper variant="outlined">
            <Typography style={{backgroundColor: "#DDD", padding: "10px"}}>{props.Meta.display_name}</Typography>
            <Button onClick={() => {
                let workCopy = []
                if (props.Data !== undefined)
                    workCopy = JSON.parse(JSON.stringify(props.Data));
                let defObject = {}
                for (let i = 0; i < props.Meta.columns.length; i++) {
                    let defValue = undefined
                    if (props.Meta.columns[i].defaultValue !== undefined) {
                        defValue = props.Meta.columns[i].defaultValue
                    } else {
                        if (props.Meta.columns[i].type === "number")
                            defValue = 0
                        if (props.Meta.columns[i].type === "text")
                            defValue = ""
                        if (props.Meta.columns[i].type === "table")
                            defValue = []
                        if (props.Meta.columns[i].type === "object")
                            defValue = {}
                    }
                    defObject[props.Meta.columns[i].name] = defValue
                }
                workCopy.push(defObject)
                props.OnChangedValue(props.Meta.name, workCopy)
            }}>Add row</Button>
            <Grid container>
                <Grid item>
                    <Table size="small" style={{minWidth: "400px"}}>
                        <TableBody>
                            <TableRow style={{backgroundColor: "#DDDDDD"}}>
                                {props.Meta.columns.map((item) => (
                                    <TableCell>
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
                                        {props.Meta.columns.map((item) => (
                                            <TableCell>
                                                <div>{(item.type !== "table" && item.type !== "object")?dataItem[item.name]:"obj"}</div>
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                )) : <div/>
                            }
                        </TableBody>
                    </Table>
                </Grid>
                <Grid item>
                    <div style={{minWidth: "100px", minHeight: "100px", margin: "5px"}}>
                        {props.Data === undefined || props.Data.constructor !== Array || currentItem === "" || currentItem >= props.Data.length || currentItem < 0?
                            <div>No row selected</div>
                            :
                            <div>
                            <ConfigObject Meta={props.Meta.columns} Data={getDataOfItem(currentItem)}
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
        </Paper>

    )
}
