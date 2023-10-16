package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/components"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func Visualize(r *Report) (*components.Page, error) {

	var p = components.NewPage()
	var err error

	if err = addGraphs(p, r.Meta, r.DataInsert, "INSERT", false); err != nil {
		return nil, fmt.Errorf("failed to generate INSERT graph: %w", err)
	}
	if err = addGraphs(p, r.Meta, r.DataGet, "GET", true); err != nil {
		return nil, fmt.Errorf("failed to generate GET graph: %w", err)
	}
	if err = addGraphs(p, r.Meta, r.DataGetAll, "GET ALL", true); err != nil {
		return nil, fmt.Errorf("failed to generate GET graph: %w", err)
	}
	if err = addGraphs(p, r.Meta, r.DataGetHeader, "GET HEADER", true); err != nil {
		return nil, fmt.Errorf("failed to generate GET graph: %w", err)
	}

	return p, nil
}

func addGraphs(p *components.Page, m Meta, v Variants, opName string, skipBytes bool) error {

	const Title = "Comparison of %s user values"
	const DescrElapse = "Chart of elapsed time for each iteration (lower is better)"
	const DescrBytes = "Chart of allocated bytes for each iteration (lower is better)"

	opName = strings.ToUpper(opName)

	var cbElapse = func(s Stamp) opts.BarData {
		s.Elapse = math.Round(s.Elapse*100) / 100
		return opts.BarData{Value: s.Elapse}
	}
	var cbBytes = func(s Stamp) opts.BarData {
		return opts.BarData{Value: s.Bytes}
	}

	var g, err = genGraph(m, v, fmt.Sprintf(Title, opName), DescrElapse, cbElapse)
	if err != nil {
		return fmt.Errorf("failed to generate elapse %s graph: %w", opName, err)
	}
	p.AddCharts(g)

	if !skipBytes {
		g, err = genGraph(m, v, fmt.Sprintf(Title, opName), DescrBytes, cbBytes)
		if err != nil {
			return fmt.Errorf("failed to generate bytes %s graph: %w", opName, err)
		}
		p.AddCharts(g)
	}

	return nil
}

func genGraph(
	m Meta, v Variants, title, descr string,
	cb func(s Stamp) opts.BarData) (components.Charter, error) {

	var bar = charts.NewBar()

	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title:    title,
			Subtitle: descr,
		}),
		charts.WithColorsOpts(opts.Colors{
			"blue", "orange", "red", "violet", "grey", "green",
		}),
		charts.WithToolboxOpts(opts.Toolbox{Show: true}),
		charts.WithLegendOpts(opts.Legend{Show: true, Top: "bottom", Orient: "horizontal"}),
	)

	//    '#c23531', '#dd6b66',
	//    '#2f4554', '#759aa0',
	//    '#61a0a8', '#e69d87',
	//    '#d48265', '#8dc1a9',
	//    '#91c7ae', '#ea7e53',
	//    '#749f83', '#eedd78',
	//    '#ca8622', '#73a373',
	//    '#bda29a', '#73b9bc',
	//    '#6e7074', '#7289ab',
	//    '#546570', '#91ca8c',
	//    '#c4ccd3', '#f49f42',

	bar.SetXAxis(m.Nums)

	for _, variant := range m.Variants {
		var res = make([]opts.BarData, 0, len(m.Nums))
		for _, num := range m.Nums {
			res = append(res, cb(v[variant][num]))
		}
		bar.AddSeries(variant, res)
	}

	bar.XYReversal()
	bar.SetSeriesOptions(charts.WithLabelOpts(opts.Label{
		Show:     true,
		Position: "right",
		Color:    "black",
	}))

	return bar, nil
}
