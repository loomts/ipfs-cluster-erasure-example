package utils

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

var diffSize = []float64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536, 131072, 262144, 524288, 1048576}
var sameSize = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}

var add_diff = []float64{0.166395887, 0.179071267, 0.17928787, 0.21515019, 0.196713244, 0.250710994, 0.203763345, 0.21878463, 0.238759477, 0.280213037, 0.386608836, 0.493297238, 1.082803653, 1.8475076449999999, 4.298925245, 5.817052986, 8.020426485, 16.342685487, 47.476508658, 99.909784545, 163.398210872}
var ecadd_diff = []float64{0.151936757, 0.147808674, 0.302904785, 0.310249541, 0.194808671, 0.189935311, 0.182074082, 0.192777818, 0.284614345, 0.349890835, 0.453092084, 0.429095897, 0.60247534, 1.322403505, 2.083115052, 5.20771203, 6.871771338, 13.402119263, 27.277070709, 63.964917357, 125.482056132}
var add_same1GB = []float64{91.184931047, 96.774071187, 98.24160554, 95.668099063, 84.502639526, 93.280970291, 94.990654178, 94.637416535, 95.898682755, 97.059757424, 99.004907566, 100.243467094, 94.992492943, 97.66777779, 100.014414768, 98.809110187, 98.057695412, 102.907488774, 103.413974143, 101.937286929, 98.571185264}
var ecadd_same1GB = []float64{66.538174167, 64.226807128, 54.910991772, 57.823912333, 55.398362483, 62.99811361, 56.507358633, 55.489491385, 58.201098864, 57.587954491, 58.101931719, 54.921586059, 67.752596118, 66.998794989, 56.746959761, 66.741377911, 66.119332201, 73.229681408, 59.574420691, 58.431667209, 56.397626786000004}

var get_diff = []float64{0.025077027, 0.004431369, 0.003879495, 0.004080458, 0.127054629, 0.003378317, 0.073612799, 0.003998472, 0.036288574, 0.085198862, 0.008046949, 0.177098612, 0.225970204, 0.047434116, 1.7756529460000001, 1.473931252, 0.351302509, 0.761663575, 2.827602645, 8.747001807, 12.046964569}
var ecget_diff = []float64{0.161326832, 0.068529851, 0.126017503, 0.044117456, 0.131322838, 0.106572067, 0.079642931, 0.035490949, 0.152357934, 0.196804474, 0.157335585, 0.080480793, 0.133246687, 1.008355691, 0.405696774, 1.770076022, 2.802202759, 6.023545095, 13.24062477, 17.261602924, 37.767332432}
var ecrecovery_get_diff = []float64{0.134461791, 0.042788602, 0.039521689, 0.043643708, 0.050496205, 0.039275368, 0.053148014, 0.04851649, 0.031158404, 0.18736684, 0.150247744, 0.21635322, 0.321617556, 0.653435399, 0.128699347, 2.508882572, 3.024677499, 1.112505926, 2.912436258, 5.944616915, 10.202298527}

var get_same = []float64{3.597560362, 17.411639056, 8.274242404, 6.157963299, 3.930695495, 16.09841406, 13.428852875, 3.710413608, 2.363980485, 3.031737762, 25.122264804, 19.920162651, 6.759335297, 4.500621443, 2.229492664, 13.45087036, 17.509822227, 15.349132354, 14.22241672, 8.126146716, 3.73750336}
var ecget_same = []float64{15.873292801, 18.531723616, 16.932024879, 14.25643669, 16.271139332, 27.355062482, 23.866623412, 15.575021985, 15.122305363, 15.898288704, 27.012587186, 20.301297035, 24.000576579, 22.287301121, 26.05079091, 23.10305274, 22.104130355, 19.937479822, 15.754206696, 18.948377859, 18.30695025}

var recon_total_time_diff = []float64{0.063925, 0.085727, 0.125465, 0.062755, 0.110150, 0.062070, 0.137806, 0.073090, 0.108912, 0.139449, 0.204580, 0.166521, 0.318681, 0.436190, 0.863551, 1.234449, 2.254924, 5.285479, 7.215680, 11.143156, 20.926259}
var recon_data_time_diff = []float64{0.016214, 0.005453, 0.043861, 0.012898, 0.010950, 0.008738, 0.052693, 0.013762, 0.056986, 0.076150, 0.116077, 0.147135, 0.157340, 0.192584, 0.304603, 0.589268, 0.968429, 1.931426, 3.278989, 5.752613, 4.425712}
var recon_repin_time_diff = []float64{0.047712, 0.080275, 0.081604, 0.049857, 0.099201, 0.053331, 0.085114, 0.059328, 0.051926, 0.063299, 0.088503, 0.019385, 0.161341, 0.243606, 0.558948, 0.645181, 1.286495, 3.354053, 3.936692, 5.390543, 16.500547}
var recon_getshards_time_diff = []float64{0.010781, 0.005406, 0.041398, 0.012152, 0.010729, 0.008502, 0.052361, 0.012378, 0.056167, 0.075855, 0.115000, 0.146053, 0.155980, 0.189918, 0.296593, 0.581127, 0.948735, 1.895303, 3.175634, 5.608609, 4.160262}
var recon_total_size_diff = []float64{3396.000000, 6468.000000, 12612.000000, 24900.000000, 49479.000000, 98631.000000, 196935.000000, 393543.000000, 786541.000000, 1048778.000000, 2097454.000000, 3670520.000000, 6816648.000000, 13108905.000000, 25693417.000000, 50862441.000000, 101200608.000000, 201876764.000000, 402739690.000000, 805476686.000000, 1610935637.000000}
var recon_data_size_diff = []float64{2264.000000, 4312.000000, 8408.000000, 24900.000000, 32986.000000, 98631.000000, 131290.000000, 393543.000000, 524288.000000, 1048576.000000, 1573166.000000, 2884088.000000, 4195208.000000, 8390313.000000, 16780521.000000, 33560937.000000, 67121888.000000, 134759436.000000, 268495832.000000, 536990264.000000, 1073954188.000000}
var recon_repin_size_diff = []float64{1132.000000, 2156.000000, 4204.000000, 0.000000, 16493.000000, 0.000000, 65645.000000, 0.000000, 262253.000000, 202.000000, 524288.000000, 786432.000000, 2621440.000000, 4718592.000000, 8912896.000000, 17301504.000000, 34078720.000000, 67117328.000000, 134243858.000000, 268486422.000000, 536981449.000000}
var recon_total_rate_diff = []float64{53124.347531, 75448.560238, 100521.847230, 396784.148929, 449195.657069, 1589038.014066, 1429069.611789, 5384385.319514, 7221814.328559, 7520878.818496, 10252498.147447, 22042414.960190, 21390175.316208, 30053206.872814, 29753200.575679, 41202532.009489, 44879830.425932, 38194600.974590, 55814512.636267, 72284431.368262, 76981540.466116}
var recon_data_rate_diff = []float64{139632.769517, 790812.055042, 191696.845392, 1930577.517921, 3012499.599304, 11287214.643957, 2491616.296632, 28596377.209455, 9200321.448304, 13769864.758343, 13552802.504797, 19601606.646552, 26663355.956000, 43567008.094474, 55089752.274263, 56953607.760386, 69310091.810210, 69771977.151879, 81883734.664800, 93347195.009206, 242662488.706813}
var recon_repin_rate_diff = []float64{23725.920182, 26857.790825, 51516.865959, 0.000000, 166259.229237, 0.000000, 771262.087311, 0.000000, 5050514.874121, 3191.213483, 5923957.781586, 40568078.554432, 16247780.134581, 19369785.835674, 15945839.575282, 26816494.124477, 26489581.042314, 20010814.934400, 34100676.780980, 49806933.313458, 32543251.005660}
var rs_recon_time_diff = []float64{0.005433, 0.000047, 0.002463, 0.000745, 0.000221, 0.000236, 0.000331, 0.001384, 0.000819, 0.000295, 0.001077, 0.001082, 0.001360, 0.002666, 0.008011, 0.008141, 0.019693, 0.036124, 0.103355, 0.144004, 0.265449}
var rs_recon_size_diff = []float64{1132.000000, 2156.000000, 4204.000000, 8300.000000, 16493.000000, 32877.000000, 65645.000000, 131181.000000, 262253.000000, 524490.000000, 1048878.000000, 2097656.000000, 4195208.000000, 8390313.000000, 16780521.000000, 33560937.000000, 67121888.000000, 134243612.000000, 268487122.000000, 536974138.000000, 1073948171.000000}
var rs_recon_rate_diff = []float64{208360.329371, 46000554.737673, 1706682.082558, 11134498.028657, 74590820.078421, 139302829.105423, 198015770.073059, 94735838.659987, 320234693.628349, 1777149052.959713, 974125553.291516, 1937432287.274672, 3083927177.350469, 3146850809.955679, 2094712650.048796, 4122229410.401313, 3408272322.641788, 3716163606.123056, 2597721407.497293, 3728872613.771161, 8169837762.048000}

//var recon_total_time_diff = []float64{}
//var recon_data_time_diff = []float64{}
//var recon_getshards_time_diff = []float64{}
//var recon_repin_time_diff = []float64{}
//var recon_total_size_diff = []float64{}
//var recon_data_size_diff = []float64{}
//var recon_repin_size_diff = []float64{}
//var recon_total_rate_diff = []float64{}
//var recon_data_rate_diff = []float64{}
//var recon_repin_rate_diff = []float64{}
//var rs_recon_time_diff = []float64{}
//var rs_recon_size_diff = []float64{}
//var rs_recon_rate_diff = []float64{}

var add_rs = []float64{3.628906, 6.628906, 12.628906, 24.628906, 48.634766, 96.634766, 192.634766, 384.634766, 768.422852, 1024.515625, 2048.780273, 3585.075195, 6657.664062, 12802.841797, 25093.185547, 49673.875000, 98835.366211, 197158.174805, 393326.108398, 786649.073242, 1573281.198242}
var add_rs_write_amplification = []float64{2.628906, 2.314453, 2.157227, 2.078613, 2.039673, 2.019836, 2.009918, 2.004959, 2.001652, 1.001007, 1.000762, 0.750525, 0.625406, 0.562847, 0.531567, 0.515926, 0.508108, 0.504198, 0.500420, 0.500414, 0.500398}
var add_3replication = []float64{3.316406, 6.316406, 12.316406, 24.316406, 48.319336, 96.319336, 192.319336, 384.319336, 768.322266, 1536.600586, 3072.893555, 6145.485352, 12290.666016, 24581.012695, 49161.700195, 98323.075195, 196646.173828, 393291.849609, 786583.382812, 1573166.437500, 3146332.614258}
var add_3replication_amplification = []float64{2.316406, 2.158203, 2.079102, 2.039551, 2.019958, 2.009979, 2.004990, 2.002495, 2.001259, 2.001173, 2.000873, 2.000725, 2.000651, 2.000612, 2.000592, 2.000582, 2.000582, 2.000579, 2.000577, 2.000577, 2.000577}

func extractECRecoveryLogInfo() {
	file, _ := os.Open("ecrecovery.log")
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ", ")
		for _, part := range parts {
			kv := strings.Split(part, ":")
			var value float64
			if strings.Contains(part, "time") {
				duration, _ := time.ParseDuration(kv[1])
				value = duration.Seconds()
			} else {
				value, _ = strconv.ParseFloat(kv[1], 64)
			}

			switch {
			case strings.Contains(part, "recon_total_time_diff"):
				recon_total_time_diff = append(recon_total_time_diff, value)
			case strings.Contains(part, "recon_getdata_time_diff"):
				recon_data_time_diff = append(recon_data_time_diff, value)
			case strings.Contains(part, "recon_repin_time_diff"):
				recon_repin_time_diff = append(recon_repin_time_diff, value)
			case strings.Contains(part, "recon_total_size_diff"):
				recon_total_size_diff = append(recon_total_size_diff, value)
			case strings.Contains(part, "recon_data_size_diff"):
				recon_data_size_diff = append(recon_data_size_diff, value)
			case strings.Contains(part, "recon_repin_size_diff"):
				recon_repin_size_diff = append(recon_repin_size_diff, value)
			case strings.Contains(part, "recon_total_rate_diff"):
				recon_total_rate_diff = append(recon_total_rate_diff, value)
			case strings.Contains(part, "recon_data_rate_diff"):
				recon_data_rate_diff = append(recon_data_rate_diff, value)
			case strings.Contains(part, "recon_repin_rate_diff"):
				recon_repin_rate_diff = append(recon_repin_rate_diff, value)
			case strings.Contains(part, "rs_recon_time_diff"):
				rs_recon_time_diff = append(rs_recon_time_diff, value)
			case strings.Contains(part, "rs_recon_size_diff"):
				rs_recon_size_diff = append(rs_recon_size_diff, value)
			case strings.Contains(part, "rs_recon_rate_diff"):
				rs_recon_rate_diff = append(rs_recon_rate_diff, value)
			}
		}
	}
	for i := 0; i < len(recon_data_time_diff); i++ {
		recon_getshards_time_diff = append(recon_getshards_time_diff, recon_data_time_diff[i]-rs_recon_time_diff[i])
	}
	printslice := func(values []float64, name string) {
		fmt.Printf("var %s = []float64{", name)
		for _, v := range values {
			fmt.Printf("%f, ", v)
		}
		fmt.Println("}")
	}
	printslice(recon_total_time_diff, "recon_total_time_diff")
	printslice(recon_data_time_diff, "recon_data_time_diff")
	printslice(recon_repin_time_diff, "recon_repin_time_diff")
	printslice(recon_getshards_time_diff, "recon_getshards_time_diff")
	printslice(recon_total_size_diff, "recon_total_size_diff")
	printslice(recon_data_size_diff, "recon_data_size_diff")
	printslice(recon_repin_size_diff, "recon_repin_size_diff")
	printslice(recon_total_rate_diff, "recon_total_rate_diff")
	printslice(recon_data_rate_diff, "recon_data_rate_diff")
	printslice(recon_repin_rate_diff, "recon_repin_rate_diff")
	printslice(rs_recon_time_diff, "rs_recon_time_diff")
	printslice(rs_recon_size_diff, "rs_recon_size_diff")
	printslice(rs_recon_rate_diff, "rs_recon_rate_diff")
}

func extraAddLogInfo() {
	file, err := os.Open("add.log")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var replication, rs []int
	var parityTotal int
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data") {
			total := parityTotal
			parityTotal = 0
			dataTotal := 0
			sizes := strings.FieldsFunc(line, func(r rune) bool { return !unicode.IsNumber(r) })
			for _, size := range sizes {
				num, err := strconv.Atoi(size)
				if err != nil {
					panic(err)
				}
				dataTotal += num
			}
			total += dataTotal
			rs = append(rs, total)
			replication = append(replication, dataTotal*3)
		} else {
			size := strings.TrimPrefix(line, "parity ")
			num, err := strconv.Atoi(size)
			if err != nil {
				panic(err)
			}
			parityTotal += num
		}
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println("rs size(KB)")
	for i := 0; i < len(rs); i++ {
		fmt.Printf("%f, ", float64(rs[i])/1024)
	}
	fmt.Println()
	fmt.Println("rs add sizerate vs origin file")
	for i := 0; i < len(rs); i++ {
		origin := 1 << i * 1024
		fmt.Printf("%f, ", float64(rs[i])/float64(origin)-1)
	}
	fmt.Println()
	fmt.Println("replication size(KB)")
	for i := 0; i < len(replication); i++ {
		fmt.Printf("%f, ", float64(replication[i])/1024)
	}
	fmt.Println()
	fmt.Println("replication add sizerate vs origin file")
	for i := 0; i < len(replication); i++ {
		origin := 1 << i * 1024
		fmt.Printf("%f, ", float64(replication[i])/float64(origin)-1)
	}
}

func extra() {
	extraAddLogInfo()
	// extractECRecoveryLogInfo()
	data := []string{
		// "2m5.38824704s",
	}
	var times []float64

	for _, d := range data {
		duration, _ := time.ParseDuration(d)
		times = append(times, duration.Seconds())
	}
	strs := make([]string, len(times))
	for i, v := range times {
		strs[i] = strconv.FormatFloat(v, 'f', -1, 64)
	}

	fmt.Println(strings.Join(strs, ", "))
}

func Draw() {
	// extra()
	DrawAddDiff()
	DrawAddSame()
	DrawGetDiff()
	DrawECRecoveryTime()
	DrawWriteAmplification()
	// Cal()
	calClusterIO()
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
	calio(1024*1024*1024*1024, 3, 0.00001, 2, 2, 2, 100*1024*1024, 1000)
	calio(2*1024*1024*1024*1024, 6, 0.00001, 4, 2, 2, 80*1024*1024, 1000)
	calio(5*1024*1024*1024*1024, 14, 0.0001, 6, 3, 3, 160*1024*1024, 1000)
	calio(6*1024*1024*1024*1024, 50, 0.001, 8, 3, 3, 200*1024*1024, 1000)
	calio(8*1024*1024*1024*1024, 500, 0.001, 5, 4, 4, 400*1024*1024, 1000)
	calio(10*1024*1024*1024*1024, 1000, 0.002, 6, 4, 4, 500*1024*1024, 1000)
}
func Cal() {
	calsame := func(times []float64, name string) {
		total := 0.0
		for _, value := range times {
			total += value
		}
		avg := float64(len(ecadd_same1GB)*512) / total
		fmt.Printf("%s rate: %.2fMB/s\n", name, avg)
	}
	caldiff := func(times []float64, sizes []float64, unit string, name string) {
		totaltime := 0.0
		for _, value := range times {
			totaltime += value
		}
		totalsize := 0.0
		for _, value := range sizes {
			totalsize += value
		}
		var avg float64
		if unit == "bit" {
			avg = totalsize / totaltime / 1024 / 1024
		} else if unit == "KB" {
			avg = totalsize / totaltime / 1024
		} else {
			fmt.Println("unit error")
		}
		fmt.Printf("%s rate: %.2fMB/s\n", name, avg)
	}
	calsame(ecadd_same1GB, "add --erasure same512")
	calsame(add_same1GB, "add --erasure same512")
	caldiff(ecadd_diff, diffSize, "KB", "add --erasure diff")
	caldiff(add_diff, diffSize, "KB", "add diff")

	calsame(get_same, "get same512")
	calsame(ecget_same, "ecget same512")
	caldiff(ecget_diff, diffSize, "KB", "ecget diff")
	caldiff(get_diff, diffSize, "KB", "get diff")
	caldiff(ecrecovery_get_diff, diffSize, "KB", "ecrecovery and get diff")

	caldiff(recon_total_time_diff, recon_total_size_diff, "bit", "ecrecovery all")
	caldiff(recon_data_time_diff, recon_data_size_diff, "bit", "ecrecovery data")
	caldiff(recon_repin_time_diff, recon_repin_size_diff, "bit", "ecrecovery repin")
	caldiff(rs_recon_time_diff, rs_recon_size_diff, "bit", "rs_recon")
}

func DrawAddDiff() {
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
	plotutil.AddLinePoints(p, "add", pointsFromXYs(diffSize, add_diff), "add --erasure", pointsFromXYs(diffSize, ecadd_diff))
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
