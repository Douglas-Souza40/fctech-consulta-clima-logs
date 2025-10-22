package service

import "regexp"

func IsValidCEP(cep string) bool {
    matched, _ := regexp.MatchString(`^\d{8}$`, cep)
    if !matched {
        return false
    }
    // Rejeita CEPs com todos os digitos iguais
    for i := 1; i < len(cep); i++ {
        if cep[i] != cep[0] {
            return true
        }
    }
    return false
}
