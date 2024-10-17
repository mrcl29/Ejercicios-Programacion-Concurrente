##################### IMPORTS #####################
import threading
import time
import random
#####################

##################### CONSTANTS GLOBALS #####################
NUM_JOGUINES = 3  # Nombre de juguetes a construir
NUM_ELFS = 6  # Nombre de elfs
NOM_ELFS = ["Taleasin", "Halafarin",
            "Ailduin", "Adamar", "Galather", "Estelar"]
CAPACITAT_SALA_ESPERA = 3  # Capacitat de la sala de espera
NUM_RENS = 9  # Nombre de rens
NOM_RENS = ["RUDOLPH", "BLITZEN", "DONDER", "CUPID",
            "COMET", "VIXEN", "PRANCER", "DANCER", "DASHER"]
#####################

##################### SEMÀFORS #####################
# Semàfor per fer esperar al Pare Noel
espera_PareNoel = threading.Semaphore(0)
# Semàfor per controlar la capacitat de la sala de espera
sala_de_espera_elfs = threading.Semaphore(CAPACITAT_SALA_ESPERA)
# Semàfor per obrir la sala de espera als nous elfs amb dubtes
espera_resolucio_dubtes = threading.Semaphore(0)
# Semàfor per controlar que s'enganxin els rens al trineu
enganxar_rens = threading.Semaphore(0)
# Semàfors que funcionen com a locks per les seccions crítiques
elf_mutex = threading.Semaphore(1)
ren_mutex = threading.Semaphore(1)
#####################

##################### VARIABLES COMPARTIDES #####################
contador_elfs_dubtes = 0  # Contador de elfs amb dubtes
contador_elfs_acabats = 0  # Contador de elfs que han acabat l'execució
contador_rens_arribats = 0  # Contador de rens que han arribat de pasturar
#####################

##################### DECLARACIÓ DE CLASSES #####################


class PareNoel(threading.Thread):  # Classe que representa al Pare Noel
    def run(self):
        global contador_elfs_dubtes, contador_rens_arribats

        print("-------> El Pare Noel diu: Estic despert però me'n torn a jeure")
        while True:
            espera_PareNoel.acquire()  # Pare Noel s'en va a dormir esperant que el despertin

            elf_mutex.acquire()  # Secció crítica
            # Verificar si és perquè hi ha elfs amb dubtes
            if contador_elfs_dubtes == CAPACITAT_SALA_ESPERA:
                print("-------> El Pare Noel diu: Atendré els dubtes d'aquests 3")
                # Espera random simulant la resolució de dubtes
                time.sleep(random.randint(0, 2))
                contador_elfs_dubtes = 0  # Reinicia el contador de elfs amb dubtes
                elf_mutex.release()  # Fi secció crítica

                for _ in range(CAPACITAT_SALA_ESPERA):
                    espera_resolucio_dubtes.release()  # Els elfs surten de la sala de espera
                print("-------> El Pare Noel diu: Estic cansat me'n torn a jeure")
            else:
                elf_mutex.release()  # Fi secció crítica

            elf_mutex.acquire()  # Secció crítica
            # Verifica si tots els elfs han acabat
            if contador_elfs_acabats == NUM_ELFS:
                elf_mutex.release()  # Fi secció crítica
                print("-------> Pare Noel diu: Les joguines estan llestes. I Els rens?")
                ren_mutex.acquire()  # Secció crítica

                while contador_rens_arribats < NUM_RENS:  # Espera a que arribin tots els rens
                    ren_mutex.release()  # Fi secció crítica
                    espera_PareNoel.acquire()
                    ren_mutex.acquire()  # Secció crítica
                ren_mutex.release()  # Fi secció crítica

                print("-------> Pare Noel diu: Enganxaré els rens i partiré")
                for _ in range(NUM_RENS):  # Enganxa cada ren
                    enganxar_rens.release()
                    time.sleep(1)  # Simula enganxar el ren

                # Acaba la execució
                print(
                    "-------> El Pare Noel ha enganxat els rens, ha carregat les joguines i se'n va")
                break
            else:
                elf_mutex.release()  # Fi secció crítica


class Elf(threading.Thread):  # Classe que representa un elf
    def __init__(self, nom):
        super().__init__()
        self.nom = nom

    def run(self):
        global contador_elfs_dubtes, contador_elfs_acabats

        # Temps de espera random inicial
        time.sleep(random.randint(1, NUM_JOGUINES * 2))
        print(f"Hola som l'elf {self.nom} construiré {NUM_JOGUINES} joguines")

        for i in range(1, NUM_JOGUINES+1):  # Per cada joguina a construir
            # Temps de espera random per simular descobrir un dubte
            time.sleep(random.uniform(0.5, 3))
            sala_de_espera_elfs.acquire()  # Entra a la sala de espera
            print(f"{self.nom} diu: tinc dubtes amb la joguina {i}")

            elf_mutex.acquire()  # Secció crítica
            contador_elfs_dubtes += 1
            if contador_elfs_dubtes == CAPACITAT_SALA_ESPERA:  # Si la sala de espera es plena avisa al Pare Noel
                print(f"{self.nom} diu: Som {
                      CAPACITAT_SALA_ESPERA} que tenim dubtes, PARE NOEEEEEL!")
                espera_PareNoel.release()  # Desperta al Pare Noel
            elf_mutex.release()  # Fi secció crítica

            espera_resolucio_dubtes.acquire()  # Espera a la resolució del dubte
            print(f"{self.nom} diu: Construeixo la joguina amb ajuda")
            time.sleep(random.uniform(0.5, 2))
            sala_de_espera_elfs.release()  # Surt de la sala de espera

        elf_mutex.acquire()  # Secció crítica
        contador_elfs_acabats += 1
        if contador_elfs_acabats == NUM_ELFS:  # Si es el darrer elf avisa al Pare Noel
            print(f"{self.nom} diu: Som el darrer avisaré al Pare Noel")
            espera_PareNoel.release()  # Desperta el Pare Noel
        elf_mutex.release()  # Fi secció crítica

        # Acaba la execució
        print(f"L'elf {self.nom} ha fet les seves joguines i acaba <---------")


class Ren(threading.Thread):  # Classe que representa un ren
    def __init__(self, nom):
        super().__init__()
        self.nom = nom

    def run(self):
        global contador_rens_arribats

        print(f"{self.nom} se'n va a pasturar")
        # Simula el temps que tarda en arribar
        time.sleep(random.randint(NUM_RENS-1, NUM_RENS*3))

        ren_mutex.acquire()  # Secció crítica
        contador_rens_arribats += 1
        if contador_rens_arribats == NUM_RENS:  # Comproba si es el darrer ren en arribar
            print(
                f"El ren {self.nom} diu: Som el darrer en voler podem partir")
            espera_PareNoel.release()  # Si es el darrer avisa al Pare Noel
        else:
            print(f"El ren {self.nom} arriba, {contador_rens_arribats}")
        ren_mutex.release()  # Fi secció crítica

        enganxar_rens.acquire()  # Espera per ser enganxat al trineu
        # Acaba la execució
        print(f"El ren {self.nom} està enganxat al trineu")


#####################


def main():
    # Comença el Pare Noel
    pareNoel = PareNoel()
    pareNoel.start()

    # Comencen els rens
    rens = [Ren(NOM_RENS[i]) for i in range(NUM_RENS)]
    for ren in rens:
        ren.start()

    # Comencen els elfs
    elfs = [Elf(NOM_ELFS[i]) for i in range(NUM_ELFS)]
    for elf in elfs:
        elf.start()

    # Esperam a que tots acabin
    for elf in elfs:
        elf.join()
    for ren in rens:
        ren.join()
    pareNoel.join()


if __name__ == "__main__":
    main()
