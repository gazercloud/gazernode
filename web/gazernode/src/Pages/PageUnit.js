import React, {useEffect, useState} from 'react';
import List from "@material-ui/core/List";
import ListItem from "@material-ui/core/ListItem";
import ListItemIcon from "@material-ui/core/ListItemIcon";
import MemoryIcon from '@material-ui/icons/Memory';
import Grid from "@material-ui/core/Grid";
import Request from "../request";
import Button from "@material-ui/core/Button";
import Typography from "@material-ui/core/Typography";

function PageUnit(props) {
    const [unitValues, setUnitValues] = React.useState([])
    const [unitState, setUnitState] = React.useState([])
    const [unitName, setUnitName] = React.useState("")

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

    const btnClick = (id, name) => {
        const nameHex = new Buffer(name).toString('hex');
        props.OnNavigate("#form=data_item&dataItemName="+ nameHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestUnitItems = (unitName) => {
        let req = {
            unit_name: unitName
        }
        Request('unit_items_values', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitValues(result)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const requestUnitState = (unitName) => {
        let req = {
            id: unitName
        }
        Request('unit_state', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnitState(result)
                            setUnitName(result.unit_name)
                            props.OnTitleUpdate("Gazer - Unit " + result.unit_name)
                            requestUnitItems(result.unit_name)
                        }
                    );
                } else {
                    res.json().then(
                        (result) => {
                            //setErrorMessage(result.error)
                        }
                    );
                }
            });
    }

    const isServiceDataItemName = (dataItemName) => {
        if (dataItemName.includes("/.service/"))
            return true
        return false
    }

    const shortDataItemName = (dataItemName) => {
        dataItemName = dataItemName.replaceAll(unitName + "/", "")
        return dataItemName
    }

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        const unitId = new Buffer(props.UnitId, 'hex').toString();
        requestUnitState(unitId)
        setFirstRendering(false)
    }

    useEffect(() => {
        const unitId = new Buffer(props.UnitId, 'hex').toString();
        const timer = setInterval(() => {
            requestUnitState(unitId)
        }, 500);
        return () => clearInterval(timer);
    });

    const displayItemName = (item, mainItem) => {
        if (mainItem) {
            return (
                <div style={{fontSize: '24pt'}}>{shortDataItemName(item.name)}</div>)
        } else {
            return (
                <div style={{fontSize: '20pt'}}>{shortDataItemName(item.name)}</div>)
        }
    }

    const formatTime = (date) => {

        let hh = date.getHours();
        if (hh < 10) hh = '0' + hh;

        let mm = date.getMinutes();
        if (mm < 10) mm = '0' + mm;

        let ss = date.getSeconds();
        if (ss < 10) ss = '0' + ss;

        let fff = date.getMilliseconds();
        if (fff < 10) {
            fff = '00' + fff;
        } else {
            if (fff < 100)
                fff = '0' + fff
        }

        return hh + ':' + mm + ':' + ss + '.' + fff;
    }

    const displayItemValue = (item, mainItem) => {
        let dt1 = new Date(item.value.t / 1000)
        let dt2 = ""
        if (dt1 !== undefined) {
            dt2 = formatTime(dt1)
        }

        let fontSize = '14pt'
        if (mainItem) {
            fontSize = '24pt'
        }

        if (item.value.u === "error") {
            return (
                <div>
                    <div style={{fontSize: fontSize, color: '#F30'}}>
                        <span>{item.value.v} </span>
                        <span>{item.value.u}</span>
                    </div>
                    <div style={{fontSize: '10pt', color: '#AAA'}}>
                        <span>{dt2}</span>
                    </div>
                </div>
            )
        } else {
            return (
                <div>
                    <div style={{fontSize: fontSize, color: '#080'}}>
                        <span>{item.value.v} </span>
                        <span>{item.value.u}</span>
                    </div>
                    <div style={{fontSize: '10pt', color: '#AAA'}}>
                        <span>{dt2}</span>
                    </div>
                </div>
            )
        }
    }

    const displayItem = (item) => {
        if (item.name === unitState.main_item) {
            return (
                <div style={{color: '#C70'}}>
                    {displayItemName(item, true)}
                    {displayItemValue(item, true)}
                </div>
            )
        } else {
            return (
                <div>
                    {displayItemName(item)}
                    {displayItemValue(item)}
                </div>
            )
        }
    }

    let name = "unknown"
    let status = "unknown"

    if (unitState !== undefined) {
        name = unitState.unit_name
        status = unitState.status
    }

    const displayItems = (onlyMain) => {
        if (unitValues === undefined || unitValues.items === undefined) {
            return (<div>loading ...</div>)
        } else {
            return (
                unitValues.items.map((item) => (
                    !isServiceDataItemName(item.name) && ((!onlyMain && unitState.main_item !== item.name) || (onlyMain && unitState.main_item === item.name)) ?
                        <ListItem
                            button
                            onMouseEnter={() => handleEnter(item.name)}
                            onMouseLeave={() => handleLeave(item.name)}
                            style={btnStyle(item.name)}
                            onClick={btnClick.bind(this, item.id, item.name)}
                        >
                            <ListItemIcon>
                                <MemoryIcon />
                            </ListItemIcon>
                            <div>
                                {displayItem(item)}
                            </div>
                        </ListItem>:<div/>
                ))
            )
        }
    }

    return (
        <div>
            <Button variant='outlined' color='primary' onClick={()=>{props.OnNavigate('#form=units')}}>
                Back to the Units
            </Button>

            <div style={{marginTop: '40px'}}>
                <Typography style={{fontSize: '10pt'}}>Unit</Typography>
                <Typography style={{fontSize: '24pt'}} color='primary'>{name}</Typography>
                <Typography style={{fontSize: '10pt'}}>Status</Typography>
                <Typography style={{fontSize: '16pt'}} color='primary'>{status}</Typography>
            </div>

            <Grid container direction="column">
                <Grid item>
                    <List component="nav">
                        {displayItems(true)}
                        {displayItems(false)}
                    </List>
                </Grid>
            </Grid>
        </div>
    );
}

export default PageUnit;
