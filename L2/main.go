package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	movies := loadMovies()
	raters := loadRatings(movies)

	res := make(chan Result, 10)
	go longestMovies(movies, res)
	go latestMovies(movies, res)
	go bestRatedMovies(movies, res)
	go mostRatedMovies(movies, res)
	res <- Result{"2. Unique Raters:", []string{fmt.Sprint(len(raters))}}

	// no in-order sync, async results.
	i, resLim := 0, 5
	for r := range res {
		if i == resLim-1 {
			close(res)
		}

		fmt.Println(r.title)
		for _, m := range r.messages {
			fmt.Println(m)
		}
		fmt.Println()

		i++
	}
}

func loadMovies() []Movie {
	movFile, err := os.Open("movies.csv")
	check(err)
	scanner := bufio.NewScanner(movFile)
	scanner.Scan() // discard first line
	movies := make([]Movie, 3200)

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		fields, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		f0, _ := strconv.ParseInt(fields[0], 10, 64)
		f2, _ := strconv.ParseInt(fields[2], 10, 64)
		f6, _ := strconv.ParseInt(fields[6], 10, 64)
		movies[i] = Movie{f0, fields[1], f2, fields[3], fields[4], fields[5], f6, fields[7], 0, 0}
	}

	return movies
}

// Updates the movies to contain the rating info,
// and returns a slice containing the raters
func loadRatings(movies []Movie) map[string]Rater {
	ratingsFile, err := os.Open("ratings.csv")
	check(err)
	scanner := bufio.NewScanner(ratingsFile)
	scanner.Scan() // discard first line
	raters := make(map[string]Rater, 1000)

	// process all the rating
	for scanner.Scan() {
		line := scanner.Text()
		reader := csv.NewReader(strings.NewReader(line))
		fields, err := reader.Read()
		if err != nil {
			log.Fatal(err)
		}

		mid, _ := strconv.ParseInt(fields[1], 10, 64)
		rating, _ := strconv.ParseFloat(fields[2], 64)

		// update the individual rater
		if rater, ok := raters[fields[0]]; ok {
			rater.avg_rating = ((float64(rater.rating_count) * rater.avg_rating) + rating) / float64(rater.rating_count+1)
			rater.rating_count += 1

			raters[fields[0]] = rater
		} else {
			nrater := Rater{rating, 1}
			raters[fields[0]] = nrater
		}

		// linear search...
		index := -1
		for i, m := range movies {
			if m.id == mid {
				index = i
			}
		}

		if index >= len(movies) || index == -1 {
			// movie not in DB
			continue
		}

		// update movie's rating record
		movies[index].avg_rating = ((float64(movies[index].rating_count) * movies[index].avg_rating) + rating) / float64(movies[index].rating_count+1)
		movies[index].rating_count += 1
	}

	return raters
}

func longestMovies(movies []Movie, results chan<- Result) {
	var longest [5]Movie
	for _, movie := range movies {
		if movie.mins >= longest[4].mins {
			idx := 0
			for movie.mins < longest[idx].mins {
				idx++ // find the index to insert at
			}
			for j := 4; j > idx; j-- {
				longest[j] = longest[j-1] // shift everything to right from end
			}
			longest[idx] = movie
		}
	}

	var messages [5]string
	for i, m := range longest {
		messages[i] = fmt.Sprint(m.title, "; Duration in mins: ", m.mins)
	}
	results <- Result{"1.a Five longest movies: ", messages[:]}
}

func latestMovies(movies []Movie, res chan<- Result) {
	var latest [5]Movie
	for _, movie := range movies {
		if movie.year > latest[4].year {
			idx := 0
			for movie.year < latest[idx].year {
				idx++ // find the index to insert at
			}
			for j := 4; j > idx; j-- {
				latest[j] = latest[j-1] // shift everything to right from end
			}
			latest[idx] = movie
		}
	}

	var messages [5]string
	for i, m := range latest {
		messages[i] = fmt.Sprint(m.title, "; Year: ", m.year)
	}

	res <- Result{"1.b Five latest movies:", messages[:]}
}

func mostRatedMovies(movies []Movie, res chan<- Result) {
	sort.Slice(movies, func(i, j int) bool {
		if movies[i].rating_count == movies[j].rating_count {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies[i].rating_count > movies[j].rating_count
	})

	var messages [5]string
	for i := 0; i < 5; i++ {
		messages[i] = fmt.Sprintf("%s; Number of Ratings: %v", movies[i].title, movies[i].rating_count)
	}

	res <- Result{"1.d Most rated movies:", messages[:]}
}

func bestRatedMovies(movies []Movie, res chan<- Result) {
	movies_filt := filterMovies(movies, func(m Movie) bool {
		return m.rating_count >= 5
	})

	// just sorting instead of doing linear time processing
	sort.Slice(movies_filt, func(i, j int) bool {
		if movies_filt[i].avg_rating == movies_filt[j].avg_rating {
			return strings.Compare(movies[i].title, movies[j].title) < 0
		}
		return movies_filt[i].avg_rating > movies_filt[j].avg_rating
	})

	var messages [5]string
	for i := 0; i < 5; i++ {
		messages[i] = fmt.Sprintf("%s; Average Rating: %v", movies_filt[i].title, movies_filt[i].avg_rating)
	}

	res <- Result{"1.c Highest rated movies:", messages[:]}
}

type Movie struct {
	id           int64
	title        string
	year         int64
	country      string
	genere       string
	director     string
	mins         int64
	poster       string
	avg_rating   float64
	rating_count int64
}

type Rater struct {
	avg_rating   float64
	rating_count int64
}

type Result struct {
	title    string
	messages []string
}

func filterMovies(m []Movie, condition func(Movie) bool) []Movie {
	filteredSlice := make([]Movie, 0)

	for _, element := range m {
		if condition(element) {
			filteredSlice = append(filteredSlice, element)
		}
	}

	return filteredSlice
}

// file io helper. https://gobyexample.com/reading-files
func check(e error) {
	if e != nil {
		panic(e)
	}
}
