package main

import (
	"log"
	"math"
	"math/rand"
	"time"
)

const (
	Q      float64 = 1     // normal factor
	T      int     = 10000 // maximum iteration period
	antNum int     = 25    // the number of ants
	u      float64 = 0.1   // incorrect ratio
	v      float64 = 0.9   // correct ratio
)

func (resources *priores) Simulate() { // start algorithm
	resources.SetDic()
	// set initial factors
	resources.answer.answer = 1e9
	var mess message
	resources.SetMess(&mess)
	switch choice {
	case 1:
		if resource0.N != 0 {
			resources.resource[0] = resource0.resource[resource0.N-1]
		} else {
			break
		}
	case 2:
		if resource1.N != 0 {
			resources.resource[0] = resource1.resource[resource1.N-1]
		} else {
			break
		}
	case 3:
		if resource2.N != 0 {
			resources.resource[0] = resource2.resource[resource2.N-1]
		} else {
			break
		}
	}
	for t := T; t != 0; t-- { // iteration period
		go resources.Iterate(&mess)
	}
}

func (resources *priores) Iterate(mess *message) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// message one ant includes
	var (
		visit   [num]bool // visit the array
		tSeq    [num]int  // path
		bestAns int       = 1e9
		bestSeq [num]int  // best solution in local
		s       = antNum
	)
	bestAns = 1e9

	s = antNum
	for ; s != 0; s-- {
		for i := 0; i < num; i++ {
			visit[i] = false
		}
		tSeq[0] = 0
		visit[0] = true

		mini := -1
		var ans float64 = 0
		for i := 1; i < resources.N; i++ { // Search the path with the maximum info factor
			temp := r.Float64()
			if temp > v { // choice incorrect
				tt := r.Int() % resources.N
				for visit[tt] {
					tt = r.Int() % resources.N
				}
				visit[tt] = true
				tSeq[i] = tt
			} else {
				ans = -1
				mini = -1
				for j := 0; j < resources.N; j++ {
					if !visit[j] && ans < mess.mess[tSeq[i-1]][j] {
						ans = mess.mess[tSeq[i-1]][j]
						mini = j
					}
				}
				visit[mini] = true
				tSeq[i] = mini
			}
		}
		//fmt.Println(tSeq[:resources.N])
		if resources.CountEnergy(tSeq) < float64(bestAns) {
			for i := 0; i < resources.N; i++ {
				bestSeq[i] = tSeq[i]
			}
			bestAns = int(resources.CountEnergy(bestSeq))
		}
	}
	if resources.CountEnergy(bestSeq) < resources.answer.answer {
		resources.seq.mux.Lock()
		for i := 0; i < resources.N; i++ {
			resources.seq.seq[i] = bestSeq[i]
		}
		resources.seq.mux.Unlock()
		resources.answer.mux.Lock()
		resources.answer.answer = resources.CountEnergy(bestSeq)
		resources.answer.mux.Unlock()
	}
	resources.Bleach(mess, bestAns)
}

func (resources *priores) Bleach(mess *message, bestAns int) { // attenuation of ant pheromone
	mess.mux.Lock()
	for i := 0; i < resources.N; i++ {
		for j := 0; j < resources.N; j++ {
			mess.mess[i][j] *= u
			mess.mess[i][j] += (1.0 - u) * Q / float64(bestAns) // attenuation algorithm
		}
	}
	mess.mux.Unlock()
}

func (resources *priores) SetDic() { // set distance
	for i := 0; i < resources.N; i++ {
		for j := 0; j < resources.N; j++ {
			resources.dic[i][j] = DicTwoPoint(resources.resource[i], resources.resource[j])
		} // set distance between resource i and resource j
	}
}

func DicTwoPoint(a, b res) float64 { // calculate the distance between two resources
	return math.Sqrt(math.Pow(float64(a.X-b.X), 2) + math.Pow(float64(a.Y-b.Y), 2))
}

func (resources *priores) SetMess(conf *message) { // set message
	conf.mux.Lock()
	for i := 0; i < resources.N; i++ {
		for j := 0; j < resources.N; j++ {
			conf.mess[i][j] = Q / resources.dic[i][j]
		}
	}
	conf.mux.Unlock()
}

func (resources *priores) CountEnergy(conf [num]int) float64 { // count the energy of one road
	var temp float64 = 0
	for i := 1; i < resources.N; i++ {
		temp += DicTwoPoint(resources.resource[conf[i]], resources.resource[conf[i-1]])
	}
	//temp += DicTwoPoint(resources.resource[conf[0]], resources.resource[conf[resources.N-1]])
	return temp
}
func (resources *priores) Init(n int, ratio float64) {
	for i := 0; i < 5; i++ {
		switch i {
		case 0:
			resource0.resource = nil
			for j := 0; j < num; j++ {
				resource0.seq.seq[j] = 0
			}
			resource0.N = 0
			resource0.answer.answer = 0
		case 1:
			resource1.resource = nil
			for j := 0; j < num; j++ {
				resource1.seq.seq[j] = 0
			}
			resource2.answer.answer = 0
		case 2:
			resource2.resource = nil
			for j := 0; j < num; j++ {
				resource2.seq.seq[j] = 0
			}
			resource2.answer.answer = 0
		case 3:
			resource3.resource = nil
			for j := 0; j < num; j++ {
				resource3.seq.seq[j] = 0
			}
			resource3.answer.answer = 0
		case 4:
			resource.resource = nil
			for j := 0; j < num; j++ {
				resource.seq.seq[j] = 0
			}
			resource.answer.answer = 0
		}
	}

	resources.N = n + 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	resources.resource = make([]res, resources.N+1)
	value := 0

	// Get data from MongoDB database
	for i := 0; i < resources.N; i++ {
		var err error
		if r.Float64() < ratio {
			resources.resource[i], err = FindData(i)
			if err != nil {
				log.Fatal(err)
			}
			value++
		}
	}
	resources.N = value

	// Get data from output.csv file
	//b := ReadCsv()
	//var (
	//	value = 0
	//	count = 1
	//)
	//for j, i := 0, 0; i < len(b); i++ {
	//	if b[i] == byte(comma) || b[i] == byte(retrn) {
	//		temp, err := strconv.ParseInt(string(b[j:i]), 0, 16)
	//		if err != nil {
	//			log.Fatal(err)
	//		}
	//		j = i + 1
	//		switch count {
	//		case 1:
	//			resources.resource[value].X = int(temp)
	//			count++
	//		case 2:
	//			resources.resource[value].Y = int(temp)
	//			count++
	//		case 3:
	//			resources.resource[value].Priority = int(temp)
	//			count = 1
	//			resources.resource[value].ID = value
	//			if r.Float64() > ratio {
	//				resources.resource[value].X = 0
	//				resources.resource[value].Y = 0
	//				resources.resource[value].Priority = 0
	//				resources.resource[value].ID = 0
	//			} else {
	//
	//				value++
	//			}
	//		}
	//	}
	//} // initialize resources' positions in .csv file
	//resources.N = value

	resources.Slice()
}

func (resources *priores) Slice() {
	var (
		zeroVar  = 0
		oneVar   = 1
		twoVar   = 1
		threeVar = 1
	)

	zeroVar, oneVar, twoVar, threeVar = 0, 1, 1, 1

	for i := 0; i < resources.N; i++ {
		switch resources.resource[i].Priority {
		case 0:
			zeroVar++
		case 1:
			oneVar++
		case 2:
			twoVar++
		case 3:
			threeVar++
		}
	}

	resource0.resource = make([]res, zeroVar)
	resource1.resource = make([]res, oneVar)
	resource2.resource = make([]res, twoVar)
	resource3.resource = make([]res, threeVar)

	zeroVar, oneVar, twoVar, threeVar = 0, 1, 1, 1
	for i := 0; i < resources.N; i++ {
		switch resources.resource[i].Priority {
		case 0:
			resource0.resource[zeroVar].X = resources.resource[i].X
			resource0.resource[zeroVar].Y = resources.resource[i].Y
			resource0.resource[zeroVar].Priority = resources.resource[i].Priority
			resource0.seq.seq[zeroVar] = i
			resource0.resource[zeroVar].ID = resource.resource[i].ID
			zeroVar++
		case 1:
			resource1.resource[oneVar].Y = resources.resource[i].Y
			resource1.resource[oneVar].X = resources.resource[i].X
			resource1.resource[oneVar].Priority = resources.resource[i].Priority
			resource1.seq.seq[oneVar] = i
			resource1.resource[oneVar].ID = resource.resource[i].ID
			oneVar++
		case 2:
			resource2.resource[twoVar].X = resources.resource[i].X
			resource2.resource[twoVar].Y = resources.resource[i].Y
			resource2.resource[twoVar].Priority = resources.resource[i].Priority
			resource2.seq.seq[twoVar] = i
			resource2.resource[twoVar].ID = resource.resource[i].ID
			twoVar++
		case 3:
			resource3.resource[threeVar].X = resources.resource[i].X
			resource3.resource[threeVar].Y = resources.resource[i].Y
			resource3.resource[threeVar].Priority = resources.resource[i].Priority
			resource3.seq.seq[threeVar] = i
			resource3.resource[threeVar].ID = resources.resource[i].ID
			threeVar++
		}
	}
	resource0.N = zeroVar
	resource1.N = oneVar
	resource2.N = twoVar
	resource3.N = threeVar
}

func (resources *priores) Output(csvChoice bool) (retrn []res, answer float64) {
	count = 0
	for i := 0; i < resource0.N; i++ {
		resources.seq.seq[count] = resource0.resource[resource0.seq.seq[i]].ID
		count++
	}
	for i := 1; i < resource1.N; i++ {
		resources.seq.seq[count] = resource1.resource[resource1.seq.seq[i]].ID
		count++
	}
	for i := 1; i < resource2.N; i++ {
		resources.seq.seq[count] = resource2.resource[resource2.seq.seq[i]].ID
		count++
	}
	for i := 1; i < resource3.N; i++ {
		resources.seq.seq[count] = resource3.resource[resource3.seq.seq[i]].ID
		count++
	}
	resources.answer.answer = resources.CountEnergy(resources.seq.seq)
	retrn = make([]res, count)
	for i := 0; i < count; i++ {
		retrn[i] = resources.resource[resources.seq.seq[i]]
		setKey(retrn[i])
	}

	if csvChoice == true {
		go resources.WriteCsv()
	}
	answer = resources.answer.answer

	// Draw a picture and output
	resources.DrawPic()
	return retrn, answer
}

/**
 * Random TSP Question with Basic Ant Colony Algorithm in Multi-Goroutine Method
 * Author: Chen Songyue
 * Create Time: 2019.5.12 22:32
 * Aim: For Intern Training by Australia & New Zealand Banking
 * Brief Introduction:
 * 	Assuming there're some resources lying in the office and the positions of them are random.
 * 	Meanwhile, The employee's position is also random.
 * 	In order to get the shortest road, we have no choice but to calculate all the possible paths, which absolutely is unreal
 * 	because it will take too much time to calculate
 * 	This kind of NP problem needs to be designed an algorithm to calculate some proper paths that can approximately become the best path
 *	 And the Ant Colony Algorithm appeared and approximately-perfectly solved such problem
 */
func realize(ratio float64, csvChoice bool) ([]res, float64) {
	resource.Init(initnumber, ratio)
	//resource.Simulate()
	for choice = 0; choice < 4; choice++ {
		switch choice {
		case 0:
			resource0.Simulate()
		case 1:
			resource1.Simulate()
		case 2:
			resource2.Simulate()
		case 3:
			resource3.Simulate()
		}
	}
	time.Sleep(time.Second)
	return resource.Output(csvChoice)
}
