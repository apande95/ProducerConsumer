package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cheggaaa/pb"
)

var mu sync.Mutex

type fileRead struct {
	word   string
	pos_x  int
	pod_id int
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//output
// var otputStr = [producerCount][100000]string{}

var producerCount int = 4
var consumerCount int = 8

const bufferSize int = 10

var searchWords []string

func produce(channel chan<- fileRead, id int, fileChnk []string, wg *sync.WaitGroup) {
	defer wg.Done()
	// words := strings.Split(fileChnk, ".")
	for i, msg := range fileChnk {
		var pa = new(fileRead)
		pa.word = msg
		pa.pos_x = i
		pa.pod_id = id
		channel <- *pa
	}
}

func consume(channel <-chan fileRead, id int, wg *sync.WaitGroup, overallBar *pb.ProgressBar, w *bufio.Writer) {
	defer wg.Done()
	for msg := range channel {
		// otputStr[msg.pod_id][msg.pos_x] = msg.word
		overallBar.Increment()
		time.Sleep(time.Millisecond)
		for _, words := range searchWords {
			if strings.Contains(msg.word, words) {
				//weirdly windows sometimes put \r (carrieage return ) instead of \n so normal replace wont do...
				re := regexp.MustCompile(`\r?\n`)
				msg.word = re.ReplaceAllString(msg.word, " ")
				// msg.word = strings.ReplaceAll(msg.word, "\n", " ")
				s := fmt.Sprintf("word \"%v\" on line %v ; %v\n", words, (msg.pod_id+1)*msg.pos_x, msg.word)
				//locking as the file is a shared resource
				mu.Lock()
				_, err := w.WriteString(s)
				check(err)
				mu.Unlock()
				w.Flush()
			}
		}
	}
}

// a rudimentry implementation of the no od lines delt with via terminal progress bar  derived from : https://pkg.go.dev/github.com/schollz/progressbar/v3#section-readme
// MIT License

// pBars := []

// func testing performance to know who took how much time

func main() {
	fmt.Println("No Of Producers: ")
	fmt.Scanln(&producerCount)

	fmt.Println("No Of Consumers: ")
	fmt.Scanln(&consumerCount)

	var wordsInp string
	fmt.Println("Words to search (, seperated ): ")
	fmt.Scanln(&wordsInp)

	searchWords = strings.Split(wordsInp, ",")

	// creating channel for buffered pipe
	channel := make(chan fileRead, bufferSize)

	// creating waitgroup for prod and consumer
	wgProd := &sync.WaitGroup{}
	wgCons := &sync.WaitGroup{}

	// initializing wait groups with the inputs
	wgProd.Add(producerCount)
	wgCons.Add(consumerCount)

	// reading the whole data dump ( to be improved by locks and such TBD )
	dat, err := os.ReadFile("sherlock.txt")
	check(err)
	datStr := string(dat)

	//output file

	t1 := time.Now().String()

	f, err := os.Create("sherlock_out_" + t1 + ".txt")
	check(err)

	defer f.Close()

	w := bufio.NewWriter(f)

	// fileLen := len(datStr)
	chunks := len(datStr) / producerCount

	// noting time
	startTime := time.Now()

	// starting prod subroutines
	var total = 0
	// bars := [consumerCount]*pb.ProgressBar{}
	for i := 0; i < producerCount; i++ {
		words := strings.Split(datStr[chunks*i:chunks*(i+1)], ".")
		// bars[i] = pb.New(len(words))
		total = total + len(words)
		go produce(channel, i, words, wgProd)
	}

	overallBar := pb.StartNew(total).Prefix("Overall ")
	// starting consumer subroutines
	for i := 0; i < consumerCount; i++ {
		go consume(channel, i, wgCons, overallBar, w)
	}

	// waiting for both ends and closing channel
	wgProd.Wait()
	close(channel)
	wgCons.Wait()
	endTime := time.Now()
	overallBar.FinishPrint("The End!")
	totalTime := endTime.Sub(startTime)
	fmt.Println("total Time: ", totalTime)

	// saving time as performace for cumulative output of all runs
	// If the file doesn't exist, create it, or append to the file
	f1, err1 := os.OpenFile("performace.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	check(err1)
	defer f1.Close()

	w1 := bufio.NewWriter(f1)
	s1 := fmt.Sprintf("%v#%v#%v\n", producerCount, consumerCount, totalTime)
	fmt.Println(s1)
	_, err2 := w1.WriteString(s1)
	check(err2)
	w1.Flush()

}
