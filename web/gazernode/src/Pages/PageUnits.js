import React, {useEffect, useState} from 'react';
import { makeStyles } from '@material-ui/core/styles';
import Grid from "@material-ui/core/Grid";
import Request from "../request";
import Typography from "@material-ui/core/Typography";
import Button from "@material-ui/core/Button";
import BlurOnIcon from "@material-ui/icons/BlurOn";

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
}));

function PageUnits(props) {
    const classes = useStyles();
    //const [sensorsList, setSensorsList] = React.useState([])
    const [units, setUnits] = React.useState({})
    //const [filterString, setFilterString] = React.useState("")

    const btnStyle = (key) => {
        let borderTop = '0px solid #333333'
        let backColor = '#1E1E1E'
        let backColorHover = '#222222'


        if (currentItem === key) {
            return {
                margin: '10px',
                borderTop: borderTop,
                backgroundColor: backColor,
            }
        } else {
            if (hoverItem === key) {
                return {
                    padding: '20px',
                    minWidth: '300px',
                    maxWidth: '300px',
                    minHeight: '200px',
                    maxHeight: '200px',
                    borderRadius: '10px',
                    margin: '10px',
                    borderTop: borderTop,
                    backgroundColor: backColor,
                }
            } else {
                return {
                    padding: '20px',
                    minWidth: '300px',
                    maxWidth: '300px',
                    minHeight: '200px',
                    maxHeight: '200px',
                    borderRadius: '10px',
                    margin: '10px',
                    borderTop: borderTop,
                    backgroundColor: backColorHover,
                }
            }
        }
    }

    const [currentItem, setCurrentItem] = useState("")
    const [hoverItem, setHoverItem] = useState("")

    const btnClick = (id, name) => {
        const idHex = new Buffer(id).toString('hex');
        props.OnNavigate("#form=unit&unitId="+ idHex)
    }

    const handleEnter = (ev, key) => {
        setHoverItem(ev)
    }

    const handleLeave = (ev, key) => {
        setHoverItem("")
    }

    const requestUnitsState = () => {
        let req = {
        }
        Request('unit_state_all', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            setUnits(result)
                            props.OnTitleUpdate('Gazer - Units')
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

    const requestStartUnit = (unitId) => {
        let req = {
            ids: [unitId]
        }
        Request('unit_start', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUnitsState()
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

    const requestStopUnit = (unitId) => {
        let req = {
            ids: [unitId]
        }
        Request('unit_stop', req)
            .then((res) => {
                if (res.status === 200) {
                    res.json().then(
                        (result) => {
                            requestUnitsState()
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

    const [firstRendering, setFirstRendering] = useState(true)
    if (firstRendering) {
        requestUnitsState()
        setFirstRendering(false)
    }

    useEffect(() => {
        const timer = setInterval(() => {
            requestUnitsState()
        }, 500);
        return () => clearInterval(timer);
    });

    const displayUnitValue = (unit) => {
        if (unit.status !== "started") {
            return (
                <div style={{
                    fontSize: '24pt',
                    color: '#555',
                    overflow: 'hidden',
                    height: '30pt',
                    width: '260px',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    textAlign: 'center'
                }}>
                    {unit.value + "  " + unit.uom}
                </div>
            )
        }

        if (unit.uom === "error") {
            return (
                <div style={{
                    fontSize: '24pt',
                    color: '#F30',
                    overflow: 'hidden',
                    height: '30pt',
                    width: '260px',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap',
                    textAlign: 'center'
                }}>
                    {unit.value + " " + unit.uom}
                </div>
            )
        }

        return (
            <div style={{
                fontSize: '24pt',
                color: '#080',
                overflow: 'hidden',
                height: '30pt',
                width: '260px',
                textOverflow: 'ellipsis',
                whiteSpace: 'nowrap',
                textAlign: 'center'
            }}>
                {unit.value + "  " + unit.uom}
            </div>
        )
    }

    const displayItem = (item) => {
        return (
            <Grid container direction='column'>
                <Grid item onClick={btnClick.bind(this, item.unit_id, item.unit_name)} style={{cursor: 'pointer'}}>
                    <Grid container>
                        <Grid item>
                            <Grid container alignItems='center'>
                                {item.status === 'started'?
                                    <BlurOnIcon color='primary' fontSize='large' style={{marginRight: '5px'}} /> :
                                    <BlurOnIcon color='disabled' fontSize='large' style={{marginRight: '5px'}} />
                                }

                                <Typography style={{
                                    fontSize: '14pt',
                                    overflow: 'hidden',
                                    height: '20pt',
                                    width: '200px',
                                    textOverflow: 'ellipsis',
                                    whiteSpace: 'nowrap',
                                }}>{item.unit_name}</Typography>
                            </Grid>
                        </Grid>
                        <Grid item style={{marginTop: '10px'}}>
                                {displayUnitValue(item)}
                        </Grid>
                    </Grid>
                </Grid>
                <Grid item style={{marginTop: '20px'}}>
                    <Button variant='outlined' color='primary' style={{minWidth: '100px', margin: '10px'}} disabled={item.status === 'started'}
                            onClick={requestStartUnit.bind(this, item.unit_id)}
                    >START</Button>
                    <Button variant='outlined' color='primary' style={{minWidth: '100px', margin: '10px'}} disabled={item.status === 'stopped'}
                            onClick={requestStopUnit.bind(this, item.unit_id)}
                    >STOP</Button>
                </Grid>
            </Grid>
        )
    }

    return (
        <div>
            <Grid container direction="column">
                <Grid item>
                    <Grid container direction="row">
                            {units !== undefined && units.items !== undefined ? units.items.map((item) => (
                                <Grid item
                                    key={"gazer-page_units-unit-" + item.unit_id}
                                    button
                                    onMouseEnter={() => handleEnter(item.unit_id)}
                                    onMouseLeave={() => handleLeave(item.unit_id)}
                                    style={btnStyle(item.unit_id)}
                                >
                                    {displayItem(item)}
                                </Grid>
                            )) : <div/>}
                        </Grid>
                </Grid>
            </Grid>
        </div>
    );
}

export default PageUnits;
