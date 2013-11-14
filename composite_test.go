package fcomposite

import(
    // "fmt"
    "testing"
    . "github.com/ricallinson/simplebdd"
)

func TestFcomposite(t *testing.T) {

    Describe("Dispatch()", func() {

        It("should return", func() {
            AssertEqual(false, true)
        })
    })

    Report(t)
}