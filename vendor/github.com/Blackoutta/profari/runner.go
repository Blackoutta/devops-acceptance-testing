package profari

import "fmt"

type Test interface {
	Run()
	Teardown()
	GetName() string
	GetErrChan() chan error
}

func RunTests(tests ...Test) (string, int) {
	rc := make(chan result)

	for _, s := range tests {
		go func(s Test) {
			r := result{
				name:    s.GetName(),
				success: true,
			}

			go func() {
				for {
					select {
					case err := <-s.GetErrChan():
						if err != nil {
							if err.Error() == "done" {
								rc <- r
								return
							}
							r.success = false
						}
					}
				}
			}()

			s.Run()
			s.Teardown()
		}(s)
	}

	var exitCode int
	template := "%-20v%-20v%-20v\n"
	resultText := fmt.Sprintf(template, "ID", "Suite", "Success?")
	for j := 0; j < len(tests); j++ {
		r := <-rc
		if r.success == false {
			exitCode = 1
		}
		r.id = j + 1
		row := fmt.Sprintf(template, r.id, r.name, r.success)
		resultText += row
	}
	return resultText, exitCode
}

type result struct {
	name    string
	success bool
	id      int
}
