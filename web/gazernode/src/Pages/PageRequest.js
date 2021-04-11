import React from 'react';
import TextField from "@material-ui/core/TextField";
import {Button, Table, TableContainer} from "@material-ui/core";
import TableBody from "@material-ui/core/TableBody";
import TableRow from "@material-ui/core/TableRow";
import TableCell from "@material-ui/core/TableCell";
import Paper from "@material-ui/core/Paper";
import Grid from "@material-ui/core/Grid";
import Typography from "@material-ui/core/Typography";

export default class PageRequest extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            error: "",
            isLoaded: false,
            items: "",
            defValue: ""
        };
    }
    handleClickRun = () => {
        this.exec(this.state.defValue)
    };

    handleClickExample = (code) => {
        this.setState({
            defValue: code,
            items: "",
            error: ""
        });
    };

    handleChange = (event) => {
        this.setState({
            error: this.state.error,
            items: this.state.items,
            defValue: event.target.value,
        });
    };

    exec(code) {
        const formData  = new FormData();
        formData.append("request", code);

        fetch("/api/request", {
            method: 'POST',
            body: formData
        })
            .then(res => res.json())
            .then(
                (result) => {
                    this.setState({
                        isLoaded: true,
                        items: JSON.stringify(result, null, 2),
                        error: ""
                    });
                },
                (error) => {
                    this.setState({
                        isLoaded: true,
                        items: "",
                        error
                    });
                }
            );
    }

    render() {
        const { error, isLoaded, items, defValue } = this.state;
        console.log("draw ...")
        return (
            <div>
                <h2>Request</h2>
                <div>
                    <Grid container direction="row" spacing={3} alignItems="flex-start">
                        <Grid item>
                            <Grid container direction="column">
                                <Grid item>
                                    <TextField
                                        id="outlined-multiline-static"
                                        label="JSON-code of request"
                                        multiline
                                        rows={25}
                                        value={defValue}
                                        defaultValue=""
                                        variant="outlined"
                                        inputProps={{
                                            style: {fontSize: 12, fontFamily: "Courier New"}
                                        }}
                                        onChange={this.handleChange.bind(this)}
                                        style={{minWidth: "500px"}}/>
                                </Grid>
                                <Grid item>
                                    <Button variant="outlined" color="primary" style={{minWidth: "500px"}} onClick={this.handleClickRun.bind(this)}>RUN</Button>
                                </Grid>
                                <Grid item>
                                    <h2>Response</h2>
                                    <pre><code>{items}</code></pre>
                                    <div>
                                        {error.error}
                                    </div>
                                </Grid>
                            </Grid>
                        </Grid>
                        <Grid item>
                            <Paper variant="outlined" style={{margin: "10px"}}>
                                <Typography style={{padding: "10px"}} variant="h6">State & Config</Typography>
                                <Grid container>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeSystemGetState)}>system_get_state</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcGetState)}>hc_get_state</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeSystemGetConfig)}>system_get_config</Button></Grid>
                                </Grid>
                            </Paper>
                            <Paper variant="outlined" style={{margin: "10px"}}>
                                <Typography style={{padding: "10px"}} variant="h6">Controlling channels</Typography>
                                <Grid container>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcSetChCode)}>hc_set_ch_code</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcSetChVoltage)}>hc_set_ch_voltage</Button></Grid>
                                </Grid>
                            </Paper>
                            <Paper variant="outlined" style={{margin: "10px"}}>
                                <Typography style={{padding: "10px"}} variant="h6">Scanning</Typography>
                                <Grid container direction="column">
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcScanStart1)}>Start - OnePage(ch0: 0-300V)</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcScanStart2)}>Start - OnePage(ch0: 0-300V, ch1: 0-150V)</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcScanStart3)}>Start - TwoPages([ch0: 0-300V, ch1: 0-150V][ch0: 300V-0, ch1: 0-300V])</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeHcScanStop)}>Stop</Button></Grid>
                                </Grid>
                            </Paper>
                            <Paper variant="outlined" style={{margin: "10px"}}>
                                <Typography style={{padding: "10px"}} variant="h6">Parameters</Typography>
                                <Grid container direction="column">
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeAddParameter)}>Add Parameter</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeUpdateParameter)}>Update parameter</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeRemoveParameter)}>Remove parameter</Button></Grid>
                                    <Grid item><Button style={{minWidth: "150px", margin: "10px"}} variant="outlined" color="primary" onClick={this.handleClickExample.bind(this, this.exampleCodeParameterSetValue)}>Set parameter value</Button></Grid>
                                </Grid>
                            </Paper>
                        </Grid>
                    </Grid>
                </div>
            </div>
        );
    }

    exampleCodeSystemGetState = `{
    "fn":"system_get_state"
}`

    exampleCodeHcGetState = `{
    "fn":"hc_get_state",
    "hcid":"hc1"
}`

    exampleCodeSystemGetConfig = `{
    "fn":"system_get_config"
}`

    exampleCodeHcSetChCode = `{
    "fn":"hc_set_ch_code",
    "hcid":"hc1",
    "ch_index": 0,
    "code": 4095
}`

    exampleCodeHcSetChVoltage = `{
    "fn":"hc_set_ch_voltage",
    "hcid":"hc1",
    "ch_index": 0,
    "voltage": 123
}`

    exampleCodeHcScanStart1 = `{
    "fn":"hc_scan_start",
    "hcid":"hc1",
    "Settings": {
         "type": 0,
         "pages": [
          {
           "repeat_count": 1,
           "step_count": 1000,
           "channels": [
            {
             "channel_index": 0,
             "ref_points": [
              {
               "step_index": 0,
               "voltage": 0
              },
              {
               "step_index": 999,
               "voltage": 300
              }
             ]
            }
           ]
          }
        ]
      }
    }`

    exampleCodeHcScanStart2 = `{
    "fn":"hc_scan_start",
    "hcid":"hc1",
    "Settings": {
         "type": 0,
         "pages": [
          {
           "repeat_count": 1,
           "step_count": 1000, "channels": [
            {
             "channel_index": 0,
             "ref_points": [
                 {"step_index": 0, "voltage": 0 },
                 {"step_index": 999, "voltage": 300}
             ]
            },
            {
             "channel_index": 1,
             "ref_points": [
                 {"step_index": 0, "voltage": 0 },
                 {"step_index": 999, "voltage": 150}
             ]
            }
           ]
          }
        ]
      }
    }`

    exampleCodeHcScanStart3 = `{
    "fn":"hc_scan_start",
    "hcid":"hc1",
    "Settings": {
         "type": 0,
         "pages": [
          {
           "repeat_count": 1,
           "step_count": 1000, "channels": [
            {
             "channel_index": 0,
             "ref_points": [
                 {"step_index": 0, "voltage": 0 },
                 {"step_index": 999, "voltage": 300}
             ]
            },
            {
             "channel_index": 1,
             "ref_points": [
                 {"step_index": 0, "voltage": 0 },
                 {"step_index": 999, "voltage": 150}
             ]
            }
           ]
          }
          ,
          {
           "repeat_count": 1,
           "step_count": 1000, "channels": [
            {
             "channel_index": 0,
             "ref_points": [
                 {"step_index": 0, "voltage": 300 },
                 {"step_index": 999, "voltage": 0}
             ]
            },
            {
             "channel_index": 1,
             "ref_points": [
                 {"step_index": 0, "voltage": 0 },
                 {"step_index": 999, "voltage": 300}
             ]
            }
           ]
          }
        ]
        
        
      }
    }`
    exampleCodeHcScanStop = `{
    "fn":"hc_scan_stop",
    "hcid":"hc1"
 }`

    exampleCodeAddParameter = `{
    "fn":"system_add_param",
    "id":"param1"
 }`

    exampleCodeUpdateParameter = `{
    "fn":"system_update_param",
    "id":"param1",
    "new_id": "param1",
    "unit_type": "hc",
    "unit_id": "hc1",
    "unit_param_id": "5"
 }`

    exampleCodeRemoveParameter = `{
    "fn":"system_remove_param",
    "id":"param1"
 }`

    exampleCodeParameterSetValue = `{
    "fn": "param_set_value",
    "p_id": "param1",
    "value": 42
 }`

}
