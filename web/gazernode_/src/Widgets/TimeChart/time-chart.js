export default function time_chart_new(elId) {
    let tc = {
        elementId_: elId,
        horScale: {},
        verticalScalesWidth: 0,
        draw: function () {
            let cnv = document.getElementById(elId)
            let ctx = cnv.getContext('2d')
            ctx.clearRect(0, 0, cnv.width, cnv.height)

            this.verticalScalesWidth = 0
            this.areas.forEach((el) => {
                if (el.verticalScaleWidth() > this.verticalScalesWidth)
                    this.verticalScalesWidth = el.verticalScaleWidth()
                el.updateVScales()
            })
            this.updateHorScale()

            this.horScale.width = cnv.width - this.verticalScalesWidth
            ctx.strokeRect(0, 0, ctx.canvas.width, ctx.canvas.height)
            let offsetY = 0
            this.areas.forEach((el) => {
                ctx.save()
                ctx.translate(0, offsetY)
                ctx.color = "black"
                el.draw(ctx)
                offsetY += el.height
                ctx.restore()
            })

            ctx.fillRect(this.mousePosX, this.mousePosY, 10, 10)

            ctx.save()
            ctx.translate(this.verticalScalesWidth, cnv.height - this.horScale.height)
            ctx.color = "black"
            this.horScale.draw(ctx)
            ctx.restore()

        },
        areas: [],
        name: "",
        addArea: function (name) {
            let area = time_chart_new_area(this)
            area.name = name
            this.areas.push(area)
            return area
        },
        updateHorScale: function () {
            let min = Number.MAX_VALUE
            let max = Number.MIN_SAFE_INTEGER

            this.areas.forEach((elArea) => {
                elArea.series.forEach((elSer) => {
                    elSer.data.forEach((elData) => {
                        if (elData[0] < min) {
                            min = elData[0]
                        }
                        if (elData[0] > max) {
                            max = elData[0]
                        }
                    })
                })
            })

            let horScaleMargin = (max - min) / 20
            this.horScale.displayMin = min - horScaleMargin
            this.horScale.displayMax = max + horScaleMargin
        },
        mousePosX: 0,
        mousePosY: 0,

        mouseMove: function (event) {
            let tc = event.target.obj
            tc.mousePosX = event.offsetX
            tc.mousePosY = event.offsetY
            tc.draw()
        },
        mouseDown: function ( event ) {
            let tc = event.target.obj
            tc.mousePosX = event.offsetX
            tc.mousePosY = event.offsetY
            tc.draw()
        },
        mouseUp: function ( event ) {
            let tc = event.target.obj
            tc.mousePosX = event.offsetX
            tc.mousePosY = event.offsetY
            tc.draw()
        },
        mouseDoubleClick: function ( event ) {
            let tc = event.target.obj
            tc.mousePosX = event.offsetX
            tc.mousePosY = event.offsetY
            tc.draw()
        }

    }

    let canvas = document.getElementById(elId)
    canvas.obj = tc
    canvas.addEventListener('mousemove', tc.mouseMove, false);
    canvas.addEventListener('mousedown', tc.mouseDown, false);
    canvas.addEventListener('mouseup', tc.mouseUp, false);
    canvas.addEventListener('mouseout', tc.mouseUp, false);
    canvas.addEventListener('dblclick', tc.mouseDoubleClick, false);

    tc.horScale = time_chart_new_horizontal_scale(tc)

    let a1 = tc.addArea("area1")
    a1.addSeries("ser 01 01")

    let a2 = tc.addArea("area2")
    a2.addSeries("ser 02 01")

    let a3 = tc.addArea("area3")
    a3.addSeries("ser 03 01")
    {
        let ser1 = a3.addSeries("ser 03 02")
        ser1.data = [[1599143745000000, 0],
            [1599143746000000, 100],
            [1599143747000000, 30],
            [1599143748000000, 40],
            [1599143749000000, 80],
            [1599143750000000, 10]]
        ser1.color = "red"
    }

    {
        let ser1 = a3.addSeries("ser 03 03")
        ser1.data = [[1599143746000000, 1100],
            [1599143747000000, 1133],
            [1599143748000000, 1130],
            [1599143749000000, 1140],
            [1599143750000000, 1180],
            [1599143751000000, 1110]]
        ser1.color = "green"
    }

    let cnv = document.getElementById(tc.elementId_)
    cnv.height = tc.areas.length * 200 + tc.horScale.height
    cnv.width = 1200

    return tc
}

function time_chart_new_area(tc) {
    return {
        tc_: tc,
        draw: function (ctx) {
            ctx.strokeRect(0, 0, ctx.canvas.width, this.height)

            this.series.forEach((el) => {
                ctx.save()
                ctx.translate(0, 0)
                el.draw(ctx)
                ctx.restore()
            })
        },
        series: [],
        name: "",
        oneVerticalScale: true,
        height: 200,
        addSeries: function (name) {
            let ser = time_chart_new_series(this)
            ser.name = name
            this.series.push(ser)
            return ser
        },
        verticalScaleWidth: function () {
            if (this.oneVerticalScale)
                return 100
            let verticalScalesWidth = 0
            this.series.forEach((el) => {
                verticalScalesWidth += el.verticalScale.width
            })
            return verticalScalesWidth
        },
        updateVScales: function() {
            if (this.oneVerticalScale) {
                let min = 1000000001
                let max = -1000000001
                this.series.forEach((el) => {
                    el.updateVScale()
                    if (el.verticalScale.displayMin < min)
                        min = el.verticalScale.displayMin
                    if (el.verticalScale.displayMax > max)
                        max = el.verticalScale.displayMax
                })
                this.series.forEach((el) => {
                    el.verticalScale.displayMin = min
                    el.verticalScale.displayMax = max
                })
            } else {
                this.series.forEach((el) => {
                    el.updateVScale()
                })
            }
        }
    }
}

function time_chart_new_series(area) {
    let newSeries = {
        area_: area,
        name: "",
        color: "blue",
        data: [[1599143744000000, 0],
            [1599143745000000, 100],
            [1599143746000000, 30],
            [1599143747000000, 40],
            [1599143748000000, 80],
            [1599143749000000, 10]],
        verticalScale: {},
        draw: function (ctx) {
            ctx.textBaseline = "top"

            ctx.save()
            let offsetX = 0
            let indexOfSeries = 0

            if (!this.area_.oneVerticalScale) {
                this.area_.series.forEach((el) => {
                    if (el === this) {
                        offsetX = indexOfSeries * 100
                    }
                    indexOfSeries++
                })
            }

            ctx.translate(offsetX, 0)
            if (this.area_.oneVerticalScale) {
                if (this.area_.series[0] === this)
                    this.verticalScale.draw(ctx)
            } else {
                this.verticalScale.draw(ctx)
            }
            ctx.restore()

            ctx.save()
            ctx.strokeStyle = this.color
            ctx.translate(this.area_.tc_.verticalScalesWidth, 0)
            ctx.beginPath()
            this.data.forEach((el) => {
                ctx.lineTo(this.area_.tc_.horScale.getPointOnX(el[0]), this.verticalScale.getPointOnY(el[1]))
            })
            ctx.stroke()
            ctx.restore()
        },
        updateVScale: function() {
            let min = 10000000000000
            let max = -10000000000000
            this.data.forEach((el) => {
                if (el[1] < min)
                    min = el[1]
                if (el[1] > max)
                    max = el[1]
            })
            let verScaleMargin = (max - min) / 20
            this.verticalScale.displayMin = min - verScaleMargin
            this.verticalScale.displayMax = max + verScaleMargin
        }
    }
    newSeries.verticalScale = time_chart_new_vertical_scale(newSeries)

    return newSeries
}

function time_chart_new_vertical_scale(ser) {
    return {
        ser_: ser,
        width: 100,
        draw: function (ctx) {
            ctx.fillStyle = "red"
            //ctx.fillRect(0, 0, 100, ser.area_.height)
            let res = this.getBeautifulScale(this.displayMin, this.displayMax, 8)

            res.forEach((el) => {
                let posInPixels = this.getPointOnY(el)
                ctx.save()
                ctx.beginPath()
                ctx.fillStyle = this.ser_.color
                ctx.strokeStyle = this.ser_.color
                ctx.textBaseline = "middle"
                ctx.textAlign = "right"
                ctx.moveTo(this.width - 10, posInPixels)
                ctx.lineTo(this.width, posInPixels)
                ctx.stroke()
                ctx.fillText(el, this.width - 15, posInPixels)
                ctx.restore()
            })

            ctx.beginPath()
            ctx.moveTo(this.width - 1, 0)
            ctx.lineTo(this.width - 1, this.ser_.area_.height)
            ctx.stroke()
        },
        series: [],
        name: "",
        height: 200,
        displayMin: 0,
        displayMax: 1,
        getBeautifulScale: function(min, max, countOfPoints) {
            if (min > max)
                return []
            if (min === max)
                return [min]

            let diapason = max - min

            // Некрасивый шаг
            let step = diapason / countOfPoints
            console.log("Step", step)

            // Порядок
            let log = Math.ceil(Math.log10(step))
            console.log(log)
            // Красивый шаг = степень 10-ки
            let step10 = Math.pow(10, log)
            console.log(step10)

            // деление на 2 - это тоже красиво
            while (diapason/(step10/2) < countOfPoints) {
                step10 = step10 / 2
            }

            // Определяем новый минимум
            let newMin = min - (min % step10)

            console.log(newMin)

            let scale = []
            // Генерируем точки
            for (let i = 0; i < countOfPoints; i++) {
                if (newMin < max && newMin > min) {
                    scale.push(newMin)
                }
                newMin += step10
            }
            return scale

        },
        getPointOnY: function (value) {
            let chartPixels = this.ser_.area_.height
            let yDelta = this.displayMax - this.displayMin
            let onePixelValue = 1
            if (Math.abs(yDelta) > 0.0001) {
                onePixelValue = chartPixels / yDelta
            }
            return chartPixels - onePixelValue * (value - this.displayMin)
        }
    }
}

function time_chart_new_horizontal_scale(tc) {
    return {
        tc_: tc,
        height: 50,
        draw: function (ctx) {
            ctx.fillStyle = "red"
            //ctx.fillRect(0, 0, 100, ser.area_.height)
            let res = this.getBeautifulScale(this.displayMin, this.displayMax, 8, 0)
            ctx.fillStyle = "red"
            //ctx.fillRect(0, 0, 100, ser.area_.height)

            res.forEach((el) => {
                let posInPixels = this.getPointOnX(el)
                ctx.save()
                ctx.beginPath()
                ctx.strokeStyle = "green"
                ctx.textBaseline = "top"
                ctx.textAlign = "center"
                ctx.moveTo(posInPixels, 0)
                ctx.lineTo(posInPixels, 10)
                ctx.stroke()

                let date = new Date(el / 1000);
                let hours = date.getHours();
                let minutes = "0" + date.getMinutes();
                let seconds = "0" + date.getSeconds();
                let formattedTime = hours + ':' + minutes.substr(-2) + ':' + seconds.substr(-2);

                ctx.fillText(formattedTime, posInPixels, 10)
                ctx.restore()
            })
        },
        series: [],
        name: "",
        width: 1000,
        displayMin: 1599143740000000,
        displayMax: 1599143750000000,
        allowedSteps: [],
        getBeautifulScale: function(min, max, countOfPoints, minStep) {
            let scale = []

            if (max < min) {
                return scale
            }

            if (max === min) {
                scale.push(min)
                return scale
            }

            ////////////////////////////////////////////////
            this.allowedSteps = []

            this.allowedSteps.push(1)      // 1 nSec
            this.allowedSteps.push(5)      // 5 nSec
            this.allowedSteps.push(10)     // 10 nSec
            this.allowedSteps.push(50)     // 50 nSec
            this.allowedSteps.push(100)    // 100 nSec
            this.allowedSteps.push(500)    // 500 nSec
            this.allowedSteps.push(1000)   // 1 mSec
            this.allowedSteps.push(5000)   // 5 mSec
            this.allowedSteps.push(10000)  // 10 mSec
            this.allowedSteps.push(50000)  // 50 mSec
            this.allowedSteps.push(100000) // 100 mSec
            this.allowedSteps.push(500000) // 500 mSec

            this.allowedSteps.push(1*1000000)  // 1 Sec
            this.allowedSteps.push(2*1000000)  // 2 Sec
            this.allowedSteps.push(5*1000000)  // 5 Sec
            this.allowedSteps.push(10*1000000) // 10 Sec
            this.allowedSteps.push(15*1000000) // 15 Sec
            this.allowedSteps.push(30*1000000) // 30 Sec

            this.allowedSteps.push(1*60*1000000)  // 1 Min
            this.allowedSteps.push(2*60*1000000)  // 2 Min
            this.allowedSteps.push(5*60*1000000)  // 5 Min
            this.allowedSteps.push(10*60*1000000) // 10 Min
            this.allowedSteps.push(15*60*1000000) // 15 Min
            this.allowedSteps.push(30*60*1000000) // 30 Min

            this.allowedSteps.push(1*60*60*1000000)  // 1 Hour
            this.allowedSteps.push(3*60*60*1000000)  // 3 Hour
            this.allowedSteps.push(6*60*60*1000000)  // 6 Hour
            this.allowedSteps.push(12*60*60*1000000) // 12 Hour

            this.allowedSteps.push(1*24*3600*1000000)    // 1 Day
            this.allowedSteps.push(2*24*3600*1000000)    // 2 Day
            this.allowedSteps.push(7*24*3600*1000000)    // 7 Day
            this.allowedSteps.push(15*24*3600*1000000)   // 15 Day
            this.allowedSteps.push(1*30*24*3600*1000000) // 1 Month
            this.allowedSteps.push(2*30*24*3600*1000000) // 2 Month
            this.allowedSteps.push(3*30*24*3600*1000000) // 3 Month
            this.allowedSteps.push(6*30*24*3600*1000000) // 3 Month
            this.allowedSteps.push(365*24*3600*1000000)  // Year
            ////////////////////////////////////////////////

            let diapason = max - min

            // Raw step - ugly
            let step = diapason / countOfPoints
            let newMin = min

            for (let i = 0; i < this.allowedSteps.length; i++) {
                let st = this.allowedSteps[i]
                if (st < minStep) {
                    continue
                }
                if (step < st) {
                    step = st // Beautiful step
                    break
                }
            }
            newMin = newMin - (newMin % step) // New begin point

            // Make points
            for (let i = 0; i < countOfPoints; i++) {
                if (newMin < max && newMin > min) {
                    scale.push(newMin)
                }
                newMin += step
            }
            return scale
        },
        getPointOnX: function (value) {
            let chartPixels = this.width
            let displayRange = this.displayMax - this.displayMin
            let offsetOfValueFromMin = value - this.displayMin
            let onePixelValue = chartPixels / displayRange
            return onePixelValue * offsetOfValueFromMin
        }
    }
}
