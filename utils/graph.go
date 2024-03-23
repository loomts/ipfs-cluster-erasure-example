package utils

import (
	"bufio"
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

//var add_diff = []float64{}
//var ecadd_diff = []float64{}
//var add_same1GB = []float64{}
//var ecadd_same1GB = []float64{}
//
//var get_diff = []float64{}
//var ecget_diff = []float64{}
//var ecrecovery_get_diff = []float64{}
//
//var get_same = []float64{}
//var ecget_same = []float64{}
//
//var recon_total_time_diff = []float64{}
//var recon_data_time_diff = []float64{}
//var recon_repin_time_diff = []float64{}
//var recon_getshards_time_diff = []float64{}
//var recon_total_size_diff = []float64{}
//var recon_data_size_diff = []float64{}
//var recon_repin_size_diff = []float64{}
//var recon_total_rate_diff = []float64{}
//var recon_data_rate_diff = []float64{}
//var recon_repin_rate_diff = []float64{}
//var rs_recon_time_diff = []float64{}
//var rs_recon_size_diff = []float64{}
//var rs_recon_rate_diff = []float64{}
//
//var add_rs = []float64{}
//var add_rs_write_amplification = []float64{}
//var add_3replication = []float64{}
//var add_3replication_amplification = []float64{}

var diffSize = []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576}
var sameSize = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

// LogEntry is a whole process of add and reconstruct a file
type LogEntry struct {
	// add
	ShardSize              int
	DataShards             int
	ParityShards           int
	RSTime                 float64
	OtherTime              float64
	TotalDataSize          int
	TotalDataAndParitySize int
	RSRate                 float64
	TotalDataRate          float64
	// recon-rs
	RsReconTimeDiff float64
	RsReconSizeDiff int
	RsReconRateDiff float64
	// recon
	ReconTotalTimeDiff   float64
	ReconGetDataTimeDiff float64
	ReconRepinTimeDiff   float64
	ReconTotalSizeDiff   int
	ReconDataSizeDiff    int
	ReconRepinSizeDiff   int
	ReconTotalRateDiff   float64
	ReconDataRateDiff    float64
	ReconRepinRateDiff   float64
	// unpin
	Name          string
	UnpinFileTime float64
}

type LogEntries []LogEntry

func (logs LogEntries) Get(name string) []float64 {
	var result []float64
	for _, log := range logs {
		switch name {
		case "ShardSize":
			result = append(result, float64(log.ShardSize))
		case "DataShards":
			result = append(result, float64(log.DataShards))
		case "ParityShards":
			result = append(result, float64(log.ParityShards))
		case "RSTime":
			result = append(result, log.RSTime)
		case "OtherTime":
			result = append(result, log.OtherTime)
		case "TotalDataSize":
			result = append(result, float64(log.TotalDataSize))
		case "TotalDataAndParitySize":
			result = append(result, float64(log.TotalDataAndParitySize))
		case "RSRate":
			result = append(result, log.RSRate)
		case "ToTalDataRate":
			result = append(result, log.TotalDataRate)
		case "rs_recon_time_diff":
			result = append(result, log.RsReconTimeDiff)
		case "recon_total_time_diff":
			result = append(result, log.ReconTotalTimeDiff)
		case "recon_getdata_time_diff":
			result = append(result, log.ReconGetDataTimeDiff)
		case "recon_repin_time_diff":
			result = append(result, log.ReconRepinTimeDiff)
		case "rs_recon_rate_diff":
			result = append(result, log.RsReconRateDiff)
		case "recon_total_rate_diff":
			result = append(result, log.ReconTotalRateDiff)
		case "recon_data_rate_diff":
			result = append(result, log.ReconDataRateDiff)
		case "recon_repin_rate_diff":
			result = append(result, log.ReconRepinRateDiff)
		case "UnpinFileTime":
			result = append(result, log.UnpinFileTime)
		}
	}
	return result
}
func Draw() {
	parseECTEST()
	// extra()
	// DrawAddDiff()
	// DrawAddSame()
	// DrawGetDiff()
	// DrawECRecoveryTime()
	// DrawWriteAmplification()
	// Cal()
	// calClusterIO()
}

func parseECTEST() {
	file, err := os.Open("addec-get.1.log")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var logs = make(LogEntries, 0)
	var log LogEntry
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", ")
		if len(parts) == 1 {
			kv := strings.Split(parts[0], " ")[1]
			k, v := strings.Split(kv, ":")[0], strings.Split(kv, ":")[1]
			log.Name = k
			value, _ := time.ParseDuration(v)
			log.UnpinFileTime = value.Seconds()
			fmt.Printf("%+v\n", log)
			logs = append(logs, log)
		}
		for _, part := range parts {
			k, v := strings.Split(part, ":")[0], strings.Split(part, ":")[1]
			switch k {
			case "ShardSize", "DataShards", "ParityShards", "TotalDataSize", "TotalDataAndParitySize", "rs_recon_size_diff", "recon_total_size_diff", "recon_data_size_diff", "recon_repin_size_diff":
				value, _ := strconv.ParseInt(v, 10, 64)
				switch k {
				case "ShardSize":
					log.ShardSize = int(value)
				case "DataShards":
					log.DataShards = int(value)
				case "ParityShards":
					log.ParityShards = int(value)
				case "TotalDataSize":
					log.TotalDataSize = int(value)
				case "TotalDataAndParitySize":
					log.TotalDataAndParitySize = int(value)
				case "rs_recon_size_diff":
					log.RsReconSizeDiff = int(value)
				case "recon_total_size_diff":
					log.ReconTotalSizeDiff = int(value)
				case "recon_data_size_diff":
					log.ReconDataSizeDiff = int(value)
				case "recon_repin_size_diff":
					log.ReconRepinSizeDiff = int(value)
				}
			case "RSTime", "OtherTime", "rs_recon_time_diff", "recon_total_time_diff", "recon_getdata_time_diff", "recon_repin_time_diff":
				duration, _ := time.ParseDuration(v)
				value := float64(duration) / float64(time.Second)
				switch k {
				case "RSTime":
					log.RSTime = value
				case "OtherTime":
					log.OtherTime = value
				case "rs_recon_time_diff":
					log.RsReconTimeDiff = value
				case "recon_total_time_diff":
					log.ReconTotalTimeDiff = value
				case "recon_getdata_time_diff":
					log.ReconGetDataTimeDiff = value
				case "recon_repin_time_diff":
					log.ReconRepinTimeDiff = value
				}
			case "RSRate", "ToTalDataRate", "rs_recon_rate_diff", "recon_total_rate_diff", "recon_data_rate_diff", "recon_repin_rate_diff":
				value, _ := strconv.ParseFloat(v, 64)
				switch k {
				case "RSRate":
					log.RSRate = value
				case "ToTalDataRate":
					log.TotalDataRate = value
				case "rs_recon_rate_diff":
					log.RsReconRateDiff = value
				case "recon_total_rate_diff":
					log.ReconTotalRateDiff = value
				case "recon_data_rate_diff":
					log.ReconDataRateDiff = value
				case "recon_repin_rate_diff":
					log.ReconRepinRateDiff = value
				}
			}
		}
	}
}

func calClusterIO() {
	//choose := func(n int, k int) float64 {
	//	if k > n {
	//		return 0
	//	}
	//	if k > n/2 {
	//		k = n - k
	//	}
	//	result := 1.0
	//	for i := 1; i <= k; i++ {
	//		result *= float64(n - k + i)
	//		result /= float64(i)
	//	}
	//	return result
	//}
	calio := func(V int, N_node int, p float64, K, M, B, S, N_add int) {
		// cal Replica_P_loss, RS_P_loss, RS-IO_cluster, Replica-IO_cluster, Delta-IO_cluster
		//$Replica\texttt{-}P_{loss} = B^{P_{down}}$
		//$Replica\texttt{-}P_{clusterloss} = 1-(1-Replica\texttt{-}P_{loss})^{N_{file}}$
		//$RS\texttt{-}P_{loss} = \sum_{i=M+1}^{K+M} {K+M \choose i} P_{down}^{i} (1-P_{down})^{K+M-i} = {k+m \choose m+1} p^{m+1} (1-p)^{k-1}$
		//$RS\texttt{-}P_{clusterloss} = 1-(1-RS\texttt{-}P_{loss})^{N_{file}}$
		//$$Replica\texttt{-}IO_{cluster} = {(B+1)}\times V_{node}\times N_{node}\times P + B\times S\times N_{add} + S\times N_{get}$$
		//$$RS\texttt{-}IO_{cluster} = \frac{(K+M)\times V_{node}\times N_{node}\times P}{K} + \frac{(K+M)\times S\times N_{add}}{K} + S\times N_{get}$$
		//$$\Delta\texttt{-}IO_{cluster} = \frac{V_{node}\times N_{node}\times P \times (K\times B-M)}{K} + S\times N_{add}\times(B-\frac{K+M}{K})$$

		//Replica_P_loss := math.Pow(p, float64(B))
		//RS_P_loss := 0.0
		//for i := M + 1; i <= K+M; i++ {
		//	RS_P_loss += math.Pow(p, float64(i)) * math.Pow(1-p, float64(K+M-i)) * choose(K+M, i)
		//}
		//Replica_clusterloss := 1 - math.Pow(1-Replica_P_loss, float64(N_file))
		//RS_clusterloss := 1 - math.Pow(1-RS_P_loss, float64(N_file))
		//Replica_storage := B
		//RS_storage := float64(K+M) / float64(K)

		Replica_IO_cluster := float64((B+1)*V*N_node)*p + float64(B*S*N_add) + float64(S*N_add)
		RS_IO_cluster := float64((K+M)*V*N_node)*p/float64(K) + float64((K+M)*S*N_add)/float64(K) + float64(S*N_add)
		detla_IO_cluster := float64(V*N_node)*p*float64(K*B-M)/float64(K) + float64(S*N_add)*(float64(B)-(float64(K+M)/float64(K)))
		fmt.Printf("\\hline\n%dTB & %d & %.5f & %d & %d & %d & %dMB & %d & %.1e & %.1e & %.1e \\\\\n", V/1024/1024/1024/1024, N_node, p, K, M, B, S/1024/1024, N_add, Replica_IO_cluster, RS_IO_cluster, detla_IO_cluster)

		//fmt.Printf("\\hline\n%d & %.5f & %d & %d & %d & %.1e & %.1e & %.1e & %.1e & %d & %.2f \\\\\n", N_file, p, K, M, B, Replica_P_loss, RS_P_loss, Replica_clusterloss, RS_clusterloss, Replica_storage, RS_storage)
	}
	calio(1024*1024*1024*1024, 4, 0.00001, 2, 2, 2, 100*1024*1024, 1000)
	calio(2*1024*1024*1024*1024, 6, 0.00001, 4, 2, 2, 80*1024*1024, 1000)
	calio(5*1024*1024*1024*1024, 9, 0.0001, 6, 3, 3, 160*1024*1024, 1000)
	calio(6*1024*1024*1024*1024, 50, 0.001, 8, 3, 3, 200*1024*1024, 1000)
	calio(8*1024*1024*1024*1024, 500, 0.001, 5, 4, 4, 400*1024*1024, 1000)
	calio(10*1024*1024*1024*1024, 5000, 0.002, 6, 4, 4, 500*1024*1024, 1000)
}

//func (logs LogEntries) Cal() {
//	calsame := func(times []float64, name string) {
//		total := 0.0
//		for _, value := range times {
//			total += value
//		}
//		avg := float64(len(ecadd_same1GB)*512) / total
//		fmt.Printf("%s rate: %.2fMB/s\n", name, avg)
//	}
//	caldiff := func(times []float64, sizes []float64, unit string, name string) {
//		totaltime := 0.0
//		for _, value := range times {
//			totaltime += value
//		}
//		totalsize := 0.0
//		for _, value := range sizes {
//			totalsize += value
//		}
//		var avg float64
//		if unit == "bit" {
//			avg = totalsize / totaltime / 1024 / 1024
//		} else if unit == "KB" {
//			avg = totalsize / totaltime / 1024
//		} else {
//			fmt.Println("unit error")
//		}
//		fmt.Printf("%s rate: %.2fMB/s\n", name, avg)
//	}
//	calsame(ecadd_same1GB, "add --erasure same512")
//	calsame(add_same1GB, "add --erasure same512")
//	caldiff(ecadd_diff, diffSize, "KB", "add --erasure diff")
//	caldiff(add_diff, diffSize, "KB", "add diff")
//
//	calsame(get_same, "get same512")
//	calsame(ecget_same, "ecget same512")
//	caldiff(ecget_diff, diffSize, "KB", "ecget diff")
//	caldiff(get_diff, diffSize, "KB", "get diff")
//	caldiff(ecrecovery_get_diff, diffSize, "KB", "ecrecovery and get diff")
//
//	caldiff(recon_total_time_diff, recon_total_size_diff, "bit", "ecrecovery all")
//	caldiff(recon_data_time_diff, recon_data_size_diff, "bit", "ecrecovery data")
//	caldiff(recon_repin_time_diff, recon_repin_size_diff, "bit", "ecrecovery repin")
//	caldiff(rs_recon_time_diff, rs_recon_size_diff, "bit", "rs_recon")
//}

func (logs LogEntries) DrawAddDiff() {
	plotter.DefaultGlyphStyle.Radius = vg.Points(3.0)
	p := plot.New()
	p.Title.Text = "add vs add --erasure (different size [1KB, 2KB ... 1GB])"
	p.X.Label.Text = "File Size(KB)"
	p.Y.Label.Text = "Time (Second)"
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = logXTicks{}
	p.Y.Tick.Marker = logYTicks{}
	plotutil.AddLinePoints(p, "add", pointsFromXYs(diffSize, logs.Get("OtherTime")), "add --erasure", pointsFromXYs(diffSize, logs.Get("OtherTime")))
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "test_add_diff.png"); err != nil {
		panic(err)
	}
}

func DrawAddSame() {
	p := plot.New()
	p.Title.Text = "add vs add --erasure (same size 1GB)"
	p.X.Label.Text = "Sequence"
	p.Y.Label.Text = "Time (Second)"
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Tick.Marker = seqTicks{}
	p.Y.Tick.Marker = seqY10Ticks{}
	p.Y.Max = 105
	plotutil.AddLinePoints(p, "add", pointsFromXYs(sameSize, add_same1GB), "add --erasure", pointsFromXYs(sameSize, ecadd_same1GB))

	if err := p.Save(12*vg.Inch, 8*vg.Inch, "test_add_same.png"); err != nil {
		panic(err)
	}
}

func DrawGetDiff() {
	p := plot.New()
	p.Title.Text = "ipfs get vs ipfs-cluster-ctl ecget (different size [1KB, 2KB ... 1GB])"
	p.X.Label.Text = "File Size(KB)"
	p.Y.Label.Text = "Time (Second)"
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.X.Scale = plot.LogScale{}
	p.X.Tick.Marker = logXTicks{}
	p.Y.Tick.Marker = logYTicks{}
	p.Y.Max = 40
	plotutil.AddLinePoints(p, "get", pointsFromXYs(diffSize, get_diff), "ecget", pointsFromXYs(diffSize, ecget_diff), "ecrecovery and get", pointsFromXYs(diffSize, ecrecovery_get_diff))
	if err := p.Save(12*vg.Inch, 8*vg.Inch, "test_get_diff.png"); err != nil {
		panic(err)
	}
}

func DrawECRecoveryTime() {
	groupA := plotter.Values(recon_getshards_time_diff)
	groupB := plotter.Values(rs_recon_time_diff)
	groupC := plotter.Values(recon_repin_time_diff)

	p := plot.New()

	p.Title.Text = "ecrecovery time (different size [1KB, 2KB ... 1GB])"
	p.X.Label.Text = "File Size(KB)"
	p.Y.Label.Text = "Time (Second)"
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Tick.Marker = seqY5Ticks{}

	w := vg.Points(10)
	barsA, err := plotter.NewBarChart(groupA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)

	barsB, err := plotter.NewBarChart(groupB, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)
	barsB.StackOn(barsA)

	barsC, err := plotter.NewBarChart(groupC, w)
	if err != nil {
		panic(err)
	}
	barsC.LineStyle.Width = vg.Length(0)
	barsC.Color = plotutil.Color(2)
	barsC.StackOn(barsB)

	p.Add(barsA, barsB, barsC)
	p.Legend.Add("get_shards_time", barsA)
	p.Legend.Add("rs_decode_time", barsB)
	p.Legend.Add("repin_time", barsC)
	p.Legend.Top = true
	p.Legend.Left = true

	str := []string{}
	for i, sz := range diffSize {
		if i == 10 {
			str = append(str, fmt.Sprintf("%.0f\n(1MB)", sz))
		} else if i == 20 {
			str = append(str, fmt.Sprintf("%.0f\n(1GB)", sz))
		} else {
			str = append(str, fmt.Sprintf("%.0f", sz))
		}
	}
	p.NominalX(str...)

	if err := p.Save(12*vg.Inch, 8*vg.Inch, "test_ecrecovery_time.png"); err != nil {
		panic(err)
	}
}

func DrawWriteAmplification() {
	groupA := plotter.Values(add_rs_write_amplification)
	groupB := plotter.Values(add_3replication_amplification)

	p := plot.New()

	p.Title.Text = "write amplification RS vs 3-Replica(different size [1KB, 2KB ... 1GB])"
	p.X.Label.Text = "File Size(KB)"
	p.Y.Label.Text = "Write Amplification Ratio(Larger than original file size)"
	p.Title.TextStyle.Font.Size = vg.Points(30)
	p.X.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Label.TextStyle.Font.Size = vg.Points(20)
	p.Y.Tick.Marker = seqY1Ticks{}

	w := vg.Points(10)
	barsA, err := plotter.NewBarChart(groupA, w)
	if err != nil {
		panic(err)
	}
	barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)

	barsB, err := plotter.NewBarChart(groupB, w)
	if err != nil {
		panic(err)
	}
	barsB.LineStyle.Width = vg.Length(0)
	barsB.Color = plotutil.Color(1)
	barsB.Offset = w

	p.Add(barsA, barsB)
	p.Legend.Add("RS", barsA)
	p.Legend.Add("3-Replica", barsB)
	p.Legend.Top = true

	str := []string{}
	for i, sz := range diffSize {
		if i == 10 {
			str = append(str, fmt.Sprintf("%.0f\n(1MB)", sz))
		} else if i == 20 {
			str = append(str, fmt.Sprintf("%.0f\n(1GB)", sz))
		} else {
			str = append(str, fmt.Sprintf("%.0f", sz))
		}
	}
	p.NominalX(str...)

	if err := p.Save(12*vg.Inch, 8*vg.Inch, "test_write_amplification.png"); err != nil {
		panic(err)
	}
}

func pointsFromXYs(x []float64, y []float64) plotter.XYs {
	pts := make(plotter.XYs, len(x))
	for i := range pts {
		pts[i].X = x[i]
		pts[i].Y = y[i]
	}
	return pts
}

type logXTicks struct{}

func (logXTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := math.Ceil(math.Log2(min)); i <= math.Floor(math.Log2(max)); i++ {
		value := math.Pow(2, i)
		var label string
		if i == 20 {
			label = fmt.Sprintf("%.0f\n(1GB)", value)
		} else if i == 10 {
			label = fmt.Sprintf("%.0f\n(1MB)", value)
		} else {
			label = fmt.Sprintf("%.0f", value)
		}
		tick := plot.Tick{Value: value, Label: label}
		ticks = append(ticks, tick)
	}
	return ticks
}

type logYTicks struct{}

func (logYTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := math.Ceil(math.Log2(min)); i <= math.Floor(math.Log2(max)); i++ {
		if i > 1 {
			break
		}
		value := math.Pow(2, i)
		if value < 1 {
			continue
		}
		var label string
		label = fmt.Sprintf("%.0f", value)
		tick := plot.Tick{Value: value, Label: label}
		ticks = append(ticks, tick)
	}
	for i := 5.0; i <= max; i *= 2 {
		label := fmt.Sprintf("%.0f", i)
		tick := plot.Tick{Value: i, Label: label}
		ticks = append(ticks, tick)
	}
	ticks = append(ticks, plot.Tick{Value: 120, Label: fmt.Sprintf("%.0f", 120.0)})
	// ticks = append(ticks, plot.Tick{Value: max, Label: fmt.Sprintf("%.0f", max)})
	return ticks
}

type seqTicks struct{}

func (seqTicks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i <= max; i++ {
		label := fmt.Sprintf("%.0f", i)
		tick := plot.Tick{Value: i, Label: label}
		ticks = append(ticks, tick)
	}
	return ticks
}

type seqY10Ticks struct{}

func (seqY10Ticks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i <= max; i += 10 {
		label := fmt.Sprintf("%.0f", i)
		tick := plot.Tick{Value: i, Label: label}
		ticks = append(ticks, tick)
	}
	return ticks
}

type seqY5Ticks struct{}

func (t seqY5Ticks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i <= max; i += 5 {
		ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.0f", i)})
	}
	ticks = append(ticks, plot.Tick{Value: 1, Label: fmt.Sprintf("%.0f", 1.0)})
	ticks = append(ticks, plot.Tick{Value: 2, Label: fmt.Sprintf("%.0f", 2.0)})
	return ticks
}

type seqY1Ticks struct{}

func (t seqY1Ticks) Ticks(min, max float64) []plot.Tick {
	var ticks []plot.Tick
	for i := min; i <= max; i += 0.5 {
		ticks = append(ticks, plot.Tick{Value: i, Label: fmt.Sprintf("%.1f", i)})
	}
	return ticks
}
