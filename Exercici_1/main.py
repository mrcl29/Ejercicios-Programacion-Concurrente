##################### IMPORTS #####################
import threading
import time
import random
#####################

##################### CONSTANTS GLOBALS #####################
NUM_JOGUINES = 3  # Número de juguetes que debe construir cada elfo
NUM_ELFS = 6  # Número de elfos
CAPACITAT_SALA_ESPERA = 3  # Capacidad de la sala de espera
NUM_RENS = 9  # Número de renos
#####################

##################### SEMÁFOROS #####################
# Semáforo para hacer esperar al Papá Noel
espera_PareNoel = threading.Semaphore(0)
# Semáforo para controlar el acceso a la sala de espera
sala_de_espera_elfs = threading.Semaphore(CAPACITAT_SALA_ESPERA)
# Semáforo para permitir el acceso de los elfos solo después de resolver sus dudas
espera_resolucion_dudas = threading.Semaphore(0)
# Semáforo para enganchar los renos al trineo
enganxar_rens = threading.Semaphore(0)
# Semáforo de contadores de elfos y renos
elf_mutex = threading.Semaphore(1)
ren_mutex = threading.Semaphore(1)
#####################

##################### VARIABLES COMPARTIDAS #####################
contador_elfs_dubtes = 0  # Contador de elfos con dudas
contador_elfs_acabats = 0  # Contador de elfos que han terminado la ejecución
contador_rens_arribats = 0  # Contador de renos que han llegado
#####################

##################### DECLARACIÓN DE CLASES #####################


class PareNoel(threading.Thread):
    def run(self):
        global contador_elfs_dubtes, contador_rens_arribats
        print("-------> El Pare Noel diu: Estic despert però me'n torn a jeure")
        while True:
            espera_PareNoel.acquire()

            # Verifica si es porque hay elfos con dudas
            elf_mutex.acquire()
            if contador_elfs_dubtes == CAPACITAT_SALA_ESPERA:
                print("-------> El Pare Noel diu: Atendré els dubtes d'aquests 3")
                # Simula la resolución de dudas
                time.sleep(random.randint(1, 2))
                contador_elfs_dubtes = 0  # Reinicia el contador de elfos con dudas
                elf_mutex.release()

                for _ in range(CAPACITAT_SALA_ESPERA):
                    espera_resolucion_dudas.release()  # Permite salir a los elfos
                print("-------> El Pare Noel diu: Estic cansat me'n torn a jeure")
            else:
                elf_mutex.release()

            # Verifica si todos los elfos han terminado
            elf_mutex.acquire()
            if contador_elfs_acabats == NUM_ELFS:
                elf_mutex.release()
                print("-------> Pare Noel diu: Les joguines estan llestes. I Els rens?")
                ren_mutex.acquire()

                # Espera a que lleguen todos los renos
                while contador_rens_arribats < NUM_RENS:
                    ren_mutex.release()
                    espera_PareNoel.acquire()
                    ren_mutex.acquire()
                ren_mutex.release()

                print("-------> Pare Noel diu: Enganxaré els rens i partiré")
                for _ in range(NUM_RENS):
                    enganxar_rens.release()
                    time.sleep(1)  # Simula enganxar el ren
                print(
                    "-------> El Pare Noel ha enganxat els rens, ha carregat les joguines i se'n va")
                break
            else:
                elf_mutex.release()


class Elf(threading.Thread):
    def __init__(self, nom):
        super().__init__()
        self.nom = nom

    def run(self):
        global contador_elfs_dubtes, contador_elfs_acabats
        print(f"Hola som l'elf {self.nom} construiré {NUM_JOGUINES} joguines")

        for i in range(1, NUM_JOGUINES+1):
            time.sleep(random.randint(0, 3))
            sala_de_espera_elfs.acquire()
            print(f"{self.nom} diu: tinc dubtes amb la joguina {i}")

            elf_mutex.acquire()
            contador_elfs_dubtes += 1
            if contador_elfs_dubtes == CAPACITAT_SALA_ESPERA:
                print(f"{self.nom} diu: Som {
                      CAPACITAT_SALA_ESPERA} que tenim dubtes, PARE NOEEEEEL!")
                espera_PareNoel.release()
            elf_mutex.release()

            espera_resolucion_dudas.acquire()  # Espera la resolución antes de continuar
            print(f"{self.nom} diu: Construeixo la joguina amb ajuda")
            time.sleep(random.randint(1, 3))
            sala_de_espera_elfs.release()

        elf_mutex.acquire()
        contador_elfs_acabats += 1
        if contador_elfs_acabats == NUM_ELFS:
            print(f"{self.nom} diu: Som el darrer avisaré al Pare Noel")
            espera_PareNoel.release()
        elf_mutex.release()

        print(f"L'elf {self.nom} ha fet les seves joguines i acaba <---------")


class Ren(threading.Thread):
    def __init__(self, nom):
        super().__init__()
        self.nom = nom

    def run(self):
        global contador_rens_arribats
        print(f"{self.nom} se'n va a pasturar")
        time.sleep(random.randint(20, 40))

        ren_mutex.acquire()
        contador_rens_arribats += 1
        if contador_rens_arribats == NUM_RENS:
            print(
                f"El ren {self.nom} diu: Som el darrer en voler podem partir")
            espera_PareNoel.release()
        else:
            print(f"El ren {self.nom} arriba, {contador_rens_arribats}")
        ren_mutex.release()

        enganxar_rens.acquire()
        print(f"El ren {self.nom} està enganxat al trineu")

#####################


def main():
    pareNoel = PareNoel()
    pareNoel.start()

    rens = [Ren(f"Ren-{i + 1}") for i in range(NUM_RENS)]
    for ren in rens:
        ren.start()

    elfs = [Elf(f"Elf-{i + 1}") for i in range(NUM_ELFS)]
    for elf in elfs:
        elf.start()

    for elf in elfs:
        elf.join()
    for ren in rens:
        ren.join()
    pareNoel.join()


if __name__ == "__main__":
    main()
